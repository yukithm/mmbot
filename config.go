package mmbot

import "fmt"

// Config for the robot.
type Config struct {
	UserName      string // Bot account name
	BindAddress   string // Bind address to listen on
	Port          int    // Port to listen on
	DisableServer bool   // Disable HTTP server
}

// Address returns bind address and port string.
func (c *Config) Address() string {
	if c.Port == 0 {
		c.Port = 8080
	}
	return fmt.Sprintf("%s:%d", c.BindAddress, c.Port)
}
