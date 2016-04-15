package app

import (
	"errors"
	"io/ioutil"
	"log"
	"mmbot"
	"mmbot/mmhook"
	"os"

	"github.com/naoina/toml"
)

// MattermostConfig is the configuration for mattermost.
type MattermostConfig struct {
	OutgoingURL        string `toml:"outgoing_url"`
	IncomingPath       string `toml:"incoming_path"`
	Token              string `toml:"token"`
	UserName           string `toml:"username"`
	OverrideUserName   string `toml:"override_username"`
	IconURL            string `toml:"icon_url"`
	InsecureSkipVerify bool   `toml:"insecure_skip_verify"`
}

// ServerConfig is the configration for the bot HTTP server.
type ServerConfig struct {
	Enable      bool   `toml:"enable"`
	BindAddress string `toml:"bind_address"`
	Port        int    `toml:"port"`
}

// Config is the configuration of the application.
type Config struct {
	Mattermost MattermostConfig `toml:"mattermost"`
	Server     ServerConfig     `toml:"server"`
}

// LoadConfigFile loads configuration file and returns Config.
func LoadConfigFile(filename string) (*Config, error) {
	var config Config
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate validates configuration values.
func (c *Config) Validate() []error {
	var errs = make([]error, 0)
	if c.Mattermost.OutgoingURL == "" {
		errs = append(errs, errors.New(`"mattermost.outgoing_url" is required`))
	}
	if c.Mattermost.UserName == "" {
		errs = append(errs, errors.New(`"mattermost.username" is required`))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ValidateAndExitOnError validates configuration values.
// Print log and exit if errors exist.
func (c *Config) ValidateAndExitOnError() {
	if errs := c.Validate(); errs != nil {
		for _, err := range errs {
			log.Printf("ERROR: %s", err)
		}
		os.Exit(1)
	}
}

// ClientConfig returns mmhook.Config.
func (c *Config) ClientConfig() *mmhook.Config {
	return &mmhook.Config{
		OutgoingURL:        c.Mattermost.OutgoingURL,
		IncomingPath:       c.Mattermost.IncomingPath,
		BindAddress:        c.Server.BindAddress,
		Port:               c.Server.Port,
		Token:              c.Mattermost.Token,
		InsecureSkipVerify: c.Mattermost.InsecureSkipVerify,
	}
}

// RobotConfig returns mmbot.Config.
func (c *Config) RobotConfig() *mmbot.Config {
	return &mmbot.Config{
		Config:           c.ClientConfig(),
		UserName:         c.Mattermost.UserName,
		OverrideUserName: c.Mattermost.OverrideUserName,
		IconURL:          c.Mattermost.IconURL,
		DisableServer:    !c.Server.Enable,
	}
}
