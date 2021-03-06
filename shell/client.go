// Package shell implements an adapter that uses readline interactive shell for development.
package shell

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/readline.v1"

	"github.com/yukithm/mmbot/adapter"
	"github.com/yukithm/mmbot/message"
	"github.com/yukithm/mmbot/mmhook"
)

// Client is a client for the shell.
type Client struct {
	config   *adapter.Config
	logger   *log.Logger
	in       chan message.InMessage
	quit     chan struct{}
	quitting bool
	errCh    chan error
}

// NewClient returns new shell client.
func NewClient(config *adapter.Config, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	c := &Client{
		config: config,
		logger: logger,
	}

	return c
}

// Start starts the interactive shell.
func (c *Client) Start() (chan message.InMessage, chan error) {
	c.in = make(chan message.InMessage, 1)
	c.quit = make(chan struct{}, 1)
	c.errCh = make(chan error, 1)

	go func() {
		c.readline()
		close(c.in)
	}()

	go func() {
		<-c.quit
		c.quitting = true
		close(c.quit)
	}()

	return c.in, c.errCh
}

// Stop terminates interactive shell.
func (c *Client) Stop() {
	c.quit <- struct{}{}
}

// Send displays a message.
func (c *Client) Send(msg *message.OutMessage) error {
	om := translateOutMessage(msg)
	if c.config.OverrideUserName != "" && om.UserName == "" {
		om.UserName = c.config.OverrideUserName
	}
	if c.config.IconURL != "" && om.IconURL == "" {
		om.IconURL = c.config.IconURL
	}

	buf, err := toJSON(om)
	if err != nil {
		return err
	}
	fmt.Printf("[Send]\n%s\n----------------\n", buf)
	fmt.Printf("mmbot> %s\n", om.Text)

	return nil
}

// IncomingWebHook returns webhook. It will be disabled if nil.
func (c *Client) IncomingWebHook() *adapter.IncomingWebHook {
	return nil
}

func (c *Client) readline() {
	rl, err := NewReadline("shell> ")
	if err != nil {
		c.errCh <- err
		return
	}
	defer rl.Close()

	for {
		time.Sleep(200 * time.Millisecond)
		line, err := rl.Readline()
		if err == io.EOF || err == readline.ErrInterrupt || c.quitting {
			return
		} else if err != nil {
			c.errCh <- err
			return
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		msg := translateInMessage(&mmhook.InMessage{
			ChannelID:   "shell",
			ChannelName: "shell",
			TeamDomain:  "shell",
			TeamID:      "shell",
			Text:        line,
			Timestamp:   strconv.FormatInt(time.Now().Unix(), 10),
			Token:       "shell_token",
			TriggerWord: strings.Fields(line)[0],
			UserID:      "shell",
			UserName:    "shell",
		})

		buf, err := toJSON(msg)
		if err != nil {
			c.errCh <- err
			return
		}
		fmt.Printf("[Receive]\n%s\n----------------\n", buf)

		c.in <- *msg
	}
}

func toJSON(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "    ")
}
