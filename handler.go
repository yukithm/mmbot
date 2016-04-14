package mmbot

import (
	"fmt"
	"regexp"
)

type Handler interface {
	CanHandle(*InMessage) bool
	Handle(*InMessage) error
}

type HandlerAction func(*InMessage) error

type PatternHandler struct {
	MessageType MessageType
	Pattern     *regexp.Regexp
	Action      HandlerAction
}

func (h *PatternHandler) CanHandle(msg *InMessage) bool {
	_, ok := h.matchPattern(msg)
	return ok
}

func (h *PatternHandler) Handle(msg *InMessage) error {
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

func (h *PatternHandler) matchPattern(msg *InMessage) ([]string, bool) {
	msgType := msg.MessageType()
	if !h.matchMessageType(msgType) {
		return nil, false
	}

	if msgType == MentionMessage {
		mentionName := msg.MentionName()
		if mentionName != msg.Robot.Config.UserName {
			return nil, false
		}
	}

	text := msg.Text
	if msgType == MentionMessage {
		text = msg.MentionlessText()
	}

	matches := h.Pattern.FindStringSubmatch(text)
	if matches == nil {
		return nil, false
	}

	return matches, true
}

func (h *PatternHandler) matchMessageType(t MessageType) bool {
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
