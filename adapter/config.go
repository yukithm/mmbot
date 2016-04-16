package adapter

// Config for an adapter.
type Config struct {
	OutgoingURL        string   // URL for incoming webhook on Mattermost
	IncomingPath       string   // Path for outgoing webhook from Mattermost
	Tokens             []string // Tokens from Mattermost
	OverrideUserName   string   // Overriding of username
	IconURL            string   // Overriding of icon URL
	InsecureSkipVerify bool     // Disable certificate checking
}
