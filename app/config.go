package app

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/naoina/toml"
	"github.com/yukithm/mmbot"
	"github.com/yukithm/mmbot/adapter"
)

// MattermostConfig is the configuration for mattermost.
type MattermostConfig struct {
	OutgoingURL        string   `toml:"outgoing_url"`
	IncomingPath       string   `toml:"incoming_path"`
	Tokens             []string `toml:"tokens"`
	UserName           string   `toml:"username"`
	OverrideUserName   string   `toml:"override_username"`
	IconURL            string   `toml:"icon_url"`
	InsecureSkipVerify bool     `toml:"insecure_skip_verify"`
}

// ServerConfig is the configration for the bot HTTP server.
type ServerConfig struct {
	Enable      bool   `toml:"enable"`
	BindAddress string `toml:"bind_address"`
	Port        int    `toml:"port"`
}

// CommonConfig is the configration of common category.
type CommonConfig struct {
	Log string `toml:"log"`
}

// Config is the configuration of the application.
type Config struct {
	Common     CommonConfig     `toml:"common"`
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

// AdapterConfig returns mmhook.Config.
func (c *Config) AdapterConfig() *adapter.Config {
	return &adapter.Config{
		OutgoingURL:        c.Mattermost.OutgoingURL,
		IncomingPath:       c.Mattermost.IncomingPath,
		Tokens:             c.Mattermost.Tokens,
		OverrideUserName:   c.Mattermost.OverrideUserName,
		IconURL:            c.Mattermost.IconURL,
		InsecureSkipVerify: c.Mattermost.InsecureSkipVerify,
	}
}

// RobotConfig returns mmbot.Config.
func (c *Config) RobotConfig() *mmbot.Config {
	return &mmbot.Config{
		UserName:      c.Mattermost.UserName,
		BindAddress:   c.Server.BindAddress,
		Port:          c.Server.Port,
		DisableServer: !c.Server.Enable,
	}
}
