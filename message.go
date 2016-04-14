package mmbot

import (
	"mmbot/mmhook"
	"regexp"
	"strings"
)

type MessageType uint

const (
	UnknownMessage MessageType = 0
	PublicMessage  MessageType = 1 << iota
	MentionMessage
	DirectMessage
	// CommandMessage
)

// InMessage represents incoming message.
type InMessage struct {
	*mmhook.InMessage          // incoming message
	Matches           []string // captured strings in the pattern
	Robot             *Robot   // robot
}

func (in *InMessage) MessageType() MessageType {
	if strings.HasPrefix(in.ChannelName, "@") {
		return DirectMessage
	}
	if strings.HasPrefix(in.Text, "@") {
		return MentionMessage
	}

	return PublicMessage
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

	msg := &mmhook.OutMessage{
		Channel: in.ChannelName,
		Text:    text,
	}
	if in.Robot.Config.OverrideUserName != "" {
		msg.UserName = in.Robot.Config.OverrideUserName
	}
	if in.Robot.Config.IconURL != "" {
		msg.IconURL = in.Robot.Config.IconURL
	}

	return in.Robot.Send(msg)
}
