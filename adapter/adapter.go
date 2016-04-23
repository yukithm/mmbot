// Package adapter defines Adapter interface for mmbot.
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
	// Start starts the communication with the chat service.
	Start() (chan message.InMessage, chan error)

	// Stop terminates the communication.
	Stop()

	// Send sends a message to chat service.
	Send(msg *message.OutMessage) error

	// IncomingWebHook returns webhook. It will be disabled if nil.
	IncomingWebHook() *IncomingWebHook
}
