package adapter

import (
	"mmbot/message"
	"net/http"
)

// IncomingWebHook represents incoming webhook that receive messages.
type IncomingWebHook struct {
	// Mount path.
	Path string

	// Handler that receive messages.
	Handler http.HandlerFunc
}

// Adapter is a client to a particular chat service.
type Adapter interface {
	// Run starts the communication with the chat service and blocks until stopped.
	Run() error

	// Stop terminates the communication.
	Stop() error

	// Send sends a message to chat service.
	Send(msg *message.OutMessage) error

	// Receiver returns a channel that receives messages from chat service.
	Receiver() chan message.InMessage

	// IncomingWebHook returns webhook. It will be disabled if nil.
	IncomingWebHook() *IncomingWebHook
}
