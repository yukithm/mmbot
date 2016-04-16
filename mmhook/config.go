package mmhook

// Config for mattermost client.
type Config struct {
	OutgoingURL        string // URL for incoming webhook on Mattermost
	IncomingPath       string // Path for outgoing webhook from Mattermost
	Token              string // Token from Mattermost
	OverrideUserName   string // Overriding of username
	IconURL            string // Overriding of icon URL
	InsecureSkipVerify bool   // Disable certificate checking
}
