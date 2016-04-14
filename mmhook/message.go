package mmhook

// InMessage represents a message from Mattermost outgouing webhook.
// (received from Mattermost)
type InMessage struct {
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	TeamDomain  string `schema:"team_domain"`
	TeamID      string `schema:"team_id"`
	Text        string `schema:"text"`
	Timestamp   string `schema:"timestamp"`
	Token       string `schema:"token"`
	TriggerWord string `schema:"trigger_word"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
}

// OutMessage represents a message to Mattermost incomig webhook.
// (send to Mattermost)
type OutMessage struct {
	Text     string `json:"text,omitempty"`
	Channel  string `json:"channel,omitempty"`
	UserName string `json:"username,omitempty"`
	IconURL  string `json:"icon_url,omitempty"`
}
