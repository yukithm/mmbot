package mmhook

import (
	"mmbot/message"
	"strings"
)

func translateInMessage(msg *InMessage) *message.InMessage {
	return &message.InMessage{
		Type:        messageType(msg),
		ChannelID:   msg.ChannelID,
		ChannelName: msg.ChannelName,
		UserID:      msg.UserID,
		UserName:    msg.UserName,
		Text:        msg.Text,
		RawMessage:  msg,
	}
}

func translateOutMessage(msg *message.OutMessage) *OutMessage {
	var channel string
	if msg.InReplyTo != nil {
		channel = msg.InReplyTo.ChannelName
	} else if msg.TriggeredBy != nil {
		channel = msg.TriggeredBy.ChannelName
	} else {
		channel = msg.ChannelName
	}

	return &OutMessage{
		Text:     msg.Text,
		Channel:  channel,
		UserName: msg.UserName,
		IconURL:  msg.IconURL,
	}
}

func messageType(msg *InMessage) message.MessageType {
	if strings.HasPrefix(msg.ChannelName, "@") {
		return message.DirectMessage
	}
	if strings.HasPrefix(msg.Text, "@") {
		return message.MentionMessage
	}

	return message.PublicMessage
}
