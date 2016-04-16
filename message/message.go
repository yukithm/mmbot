package message

import (
	"regexp"
	"strings"
)

type Sender interface {
	Send(*OutMessage) error
	SenderName() string
}

type Type uint

const (
	UnknownMessage Type = 0
	PublicMessage  Type = 1 << iota
	MentionMessage
	DirectMessage
	// CommandMessage
)

// InMessage represents an incoming message.
type InMessage struct {
	Sender      Sender
	Matches     []string // captured strings in the pattern
	Type        Type
	ChannelID   string
	ChannelName string
	UserID      string
	UserName    string
	Text        string
	RawMessage  interface{} // adapter's raw message data
}

// OutMessage represents an outgoing message.
type OutMessage struct {
	ChannelID   string
	ChannelName string
	UserName    string
	IconURL     string
	Text        string
	InReplyTo   *InMessage // reply target message
	TriggeredBy *InMessage // trigger source message
}

var mentionNameRegexp = regexp.MustCompile(`\A@([0-9a-zA-Z_]+)`)
var mentionPrefixRegexp = regexp.MustCompile(`\A@(?:[0-9a-zA-Z_]+)\s*`)

func (in *InMessage) MentionName() string {
	if matches := mentionNameRegexp.FindStringSubmatch(in.Text); matches != nil {
		return matches[1]
	}

	return ""
}

func (in *InMessage) MentionlessText() string {
	if loc := mentionPrefixRegexp.FindStringIndex(in.Text); loc != nil {
		return in.Text[loc[1]:]
	}

	return in.Text
}

// Reply sends a reply message to the sender.
func (in *InMessage) Reply(text string) error {
	targetUser := "@" + in.UserName
	if !strings.HasPrefix(text, targetUser) {
		text = targetUser + " " + text
	}

	msg := &OutMessage{
		ChannelID:   in.ChannelID,
		ChannelName: in.ChannelName,
		Text:        text,
		InReplyTo:   in,
		TriggeredBy: in,
	}

	return in.Sender.Send(msg)
}
