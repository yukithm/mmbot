package app

import (
	"errors"
	"io/ioutil"
	"log"
	"mmbot"
	"mmbot/mmhook"
	"os"

	"gopkg.in/yaml.v2"
)

// MattermostConfig is the configuration for mattermost.
type MattermostConfig struct {
	OutgoingURL        string `yaml:"outgoing_url"`
	IncomingPath       string `yaml:"incoming_path"`
	Token              string `yaml:"token"`
	UserName           string `yaml:"username"`
	OverrideUserName   string `yaml:"override_username"`
	IconURL            string `yaml:"icon_url"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

// ServerConfig is the configration for the bot HTTP server.
type ServerConfig struct {
	Enable      bool   `yaml:"enable"`
	BindAddress string `yaml:"bind_address"`
	Port        int    `yaml:port`
}

// Config is the configuration of the application.
type Config struct {
	Mattermost MattermostConfig
	Server     ServerConfig
}

func DefaultConfig() *Config {
	return &Config{
		Mattermost: MattermostConfig{
			IncomingPath: "/",
		},
		Server: ServerConfig{
			Enable:      true,
			BindAddress: "",
			Port:        8080,
		},
	}
}

// LoadConfigFile loads configuration file and returns Config.
func LoadConfigFile(filename string) (*Config, error) {
	config := DefaultConfig()
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
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
