package mmhook

import "fmt"

// Config for mattermost client.
type Config struct {
	URL                string // URL for incoming webhook on Mattermost
	IncomingURL        string // URL for outgoing webhook from Mattermost
	BindAddress        string // Bind address to listen on
	Port               int    // Port to listen on
	Token              string // Token from Mattermost
	InsecureSkipVerify bool   // Disable certificate checking
}

// Address returns bind address and port string.
func (c *Config) Address() string {
	if c.Port == 0 {
		c.Port = 8080
	}
	return fmt.Sprintf("%s:%d", c.BindAddress, c.Port)
}
