package shell

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mmbot/adapter"
	"mmbot/message"
	"mmbot/mmhook"
	"strconv"
	"strings"
	"time"

	"gopkg.in/readline.v1"
)

// Client is a client for the shell.
type Client struct {
	config *adapter.Config
	logger *log.Logger
	in     chan message.InMessage
	quit   chan bool
}

// NewClient returns new shell client.
func NewClient(config *adapter.Config, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	c := &Client{
		config: config,
		logger: logger,
		in:     make(chan message.InMessage),
		quit:   make(chan bool),
	}

	return c
}

// Run starts the interactive shell and blocks until stopped.
func (c *Client) Run() error {
	go c.readline()
	<-c.quit
	close(c.quit)
	return nil
}

// Stop terminates interactive shell.
func (c *Client) Stop() error {
	c.quit <- true
	return nil
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

// Receiver returns a channel that receives messages from shell.
func (c *Client) Receiver() chan message.InMessage {
	return c.in
}

// IncomingWebHook returns webhook. It will be disabled if nil.
func (c *Client) IncomingWebHook() *adapter.IncomingWebHook {
	return nil
}

func (c *Client) readline() {
	rl, err := NewReadline("shell> ")
	if err != nil {
		c.logger.Println(err)
		close(c.in)
		return
	}
	defer rl.Close()

	for {
		time.Sleep(200 * time.Millisecond)
		line, err := rl.Readline()
		if err == io.EOF || err == readline.ErrInterrupt {
			close(c.in)
			return
		} else if err != nil {
			c.logger.Println(err)
			close(c.in)
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
			c.logger.Println(err)
			close(c.in)
			return
		}
		fmt.Printf("[Receive]\n%s\n----------------\n", buf)

		c.in <- *msg
	}
}

func toJSON(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "    ")
}
