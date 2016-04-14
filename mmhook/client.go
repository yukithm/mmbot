package mmhook

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

// Client is a client for Mattermost.
type Client struct {
	*Config
	http *http.Client
	In   chan InMessage
}

// NewClient returns new mattermost client.
func NewClient(config *Config) *Client {
	c := &Client{
		Config: config,
		In:     make(chan InMessage),
	}
	if config.InsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
		}
		c.http = &http.Client{Transport: tr}
	} else {
		c.http = &http.Client{}
	}

	return c
}

// Send sends a message to Mattermost.
func (c *Client) Send(msg *OutMessage) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	res, err := c.http.Post(c.OutgoingURL, "application/json", bytes.NewReader(buf))
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

// Receive receives a message from Mattermost.
func (c *Client) Receive() *InMessage {
	msg := <-c.In
	return &msg
}

// StartServer runs HTTP server for incoming messages.
// (from Mattermost outgoing webhook)
func (c *Client) StartServer() {
	go c.startServer()
}

func (c *Client) startServer() {
	path := c.IncomingPath
	if path == "" {
		path = "/"
	}
	mux := http.NewServeMux()
	mux.Handle(path, c)
	log.Printf("Listening on %s\n", c.Address())
	if err := http.ListenAndServe(c.Address(), mux); err != nil {
		log.Fatal(err)
	}
}

// ServeHTTP implements http.Handler interface.
// ServeHTTP receives a message from Mattermost outgoing webhook.
func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("Invalid %q request from %q", r.Method, r.RemoteAddr)
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	msg := InMessage{}
	if err := decodeForm(&msg, r); err != nil {
		log.Printf("Invalid form data: %v", err)
		http.NotFound(w, r)
		return
	}

	if c.Token != "" {
		if msg.Token == "" {
			log.Printf("No token request from %q", r.RemoteAddr)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		} else if msg.Token != c.Token {
			log.Printf("Invalid token %q request from %q", msg.Token, r.RemoteAddr)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	c.In <- msg
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
