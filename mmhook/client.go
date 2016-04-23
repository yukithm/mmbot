package mmhook

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mmbot/adapter"
	"mmbot/message"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

// Client is a client for Mattermost.
type Client struct {
	config *adapter.Config
	logger *log.Logger
	http   *http.Client
	tokens map[string]int
	in     chan message.InMessage
	errCh  chan error
}

// NewClient returns new mattermost webhook client.
func NewClient(config *adapter.Config, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	c := &Client{
		config: config,
		logger: logger,
	}

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if config.InsecureSkipVerify {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		}
	}
	c.http = &http.Client{Transport: tr}

	// build token lookup table
	c.tokens = make(map[string]int, len(config.Tokens))
	for i, token := range config.Tokens {
		token = strings.TrimSpace(token)
		if token != "" {
			c.tokens[token] = i
		}
	}

	return c
}

// Start starts the communication with Mattermost.
func (c *Client) Start() (chan message.InMessage, chan error) {
	c.in = make(chan message.InMessage, 1)
	c.errCh = make(chan error, 1)
	return c.in, c.errCh
}

// Stop terminates the communication.
func (c *Client) Stop() {
	close(c.in)
	close(c.errCh)
}

// Send sends a message to Mattermost.
func (c *Client) Send(msg *message.OutMessage) error {
	om := translateOutMessage(msg)
	if c.config.OverrideUserName != "" && om.UserName == "" {
		om.UserName = c.config.OverrideUserName
	}
	if c.config.IconURL != "" && om.IconURL == "" {
		om.IconURL = c.config.IconURL
	}

	buf, err := json.Marshal(om)
	if err != nil {
		return err
	}

	res, err := c.http.Post(c.config.OutgoingURL, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	io.Copy(ioutil.Discard, res.Body)
	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to send a message (%d %s)",
			res.StatusCode, res.Status)
	}

	return nil
}

// IncomingWebHook returns webhook. It will be disabled if nil.
func (c *Client) IncomingWebHook() *adapter.IncomingWebHook {
	return &adapter.IncomingWebHook{
		Path:    c.config.IncomingPath,
		Handler: c.ServeHTTP,
	}
}

// ServeHTTP implements http.Handler interface.
// ServeHTTP receives a message from Mattermost outgoing webhook.
func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		c.logger.Printf("Invalid %q request from %q", r.Method, r.RemoteAddr)
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	msg := InMessage{}
	if err := decodeForm(&msg, r); err != nil {
		c.logger.Printf("Invalid form data: %v", err)
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if len(c.tokens) > 0 {
		if msg.Token == "" {
			c.logger.Printf("No token request from %q", r.RemoteAddr)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		} else if !c.validToken(msg.Token) {
			c.logger.Printf("Invalid token %q request from %q", msg.Token, r.RemoteAddr)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}
	}

	im := translateInMessage(&msg)
	c.in <- *im
}

func (c *Client) validToken(token string) bool {
	_, ok := c.tokens[token]
	return ok
}

func decodeForm(msg *InMessage, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	decoder := schema.NewDecoder()
	if err := decoder.Decode(msg, r.PostForm); err != nil {
		return err
	}

	return nil
}
