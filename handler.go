package mmbot

import (
	"fmt"
	"mmbot/message"
	"regexp"
)

type Handler interface {
	CanHandle(*message.InMessage) bool
	Handle(*message.InMessage) error
}

type HandlerAction func(*message.InMessage) error

type PatternHandler struct {
	MessageType message.MessageType
	Pattern     *regexp.Regexp
	Action      HandlerAction
}

func (h *PatternHandler) CanHandle(msg *message.InMessage) bool {
	_, ok := h.matchPattern(msg)
	return ok
}

func (h *PatternHandler) Handle(msg *message.InMessage) error {
	matches, ok := h.matchPattern(msg)
	if !ok {
		return fmt.Errorf("Cannot handle message: %#v", msg)
	}
	msg.Matches = matches

	if err := h.Action(msg); err != nil {
		return err
	}

	return nil
}

func (h *PatternHandler) matchPattern(msg *message.InMessage) ([]string, bool) {
	if !h.matchMessageType(msg.Type) {
		return nil, false
	}

	if msg.Type == message.MentionMessage {
		mentionName := msg.MentionName()
		if mentionName != msg.Sender.SenderName() {
			return nil, false
		}
	}

	text := msg.Text
	if msg.Type == message.MentionMessage {
		text = msg.MentionlessText()
	}

	matches := h.Pattern.FindStringSubmatch(text)
	if matches == nil {
		return nil, false
	}

	return matches, true
}

func (h *PatternHandler) matchMessageType(t message.MessageType) bool {
	if h.MessageType == 0 {
		return true
	}
	return t&h.MessageType != 0
}

func (h *PatternHandler) trimBotName(text string, name string) string {
	pattern := fmt.Sprintf(`\A@?(?:%s)\s*[:,]?\s+`, regexp.QuoteMeta(name))
	re := regexp.MustCompile(pattern)
	if loc := re.FindStringIndex(text); loc != nil {
		return text[loc[1]:]
	}

	return text
}
