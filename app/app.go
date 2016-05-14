// Package app provides a base of the bot application.
// App handles command line arguments and manages mmbot.Robot.
package app

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/kardianos/osext"
	"github.com/mitchellh/go-homedir"
	"github.com/yukithm/mmbot"
)

// App is a bot application.
type App struct {
	*cli.App
	Config       *Config
	ConfigLoader func(file string) (*Config, error)
	InitRobot    func(*mmbot.Robot)
}

// NewApp creates new App.
func NewApp() *App {
	app := &App{
		App: cli.NewApp(),
	}

	// app.App.Name = "mmbot"
	app.App.Usage = "a bot for Mattermost"

	app.App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Usage: "configuration file",
		},
	}

	app.App.Commands = []cli.Command{
		app.newNewConfigCommand(),
		app.newRunCommand(),
		app.newShellCommand(),
	}

	return app
}

// LoadConfig loads configuration file and store it into app.Config.
func (app *App) LoadConfig(c *cli.Context) error {
	file := getContextConfigFile(c)
	if app.ConfigLoader == nil {
		app.ConfigLoader = defaultConfigLoader
	}

	config, err := app.ConfigLoader(file)
	if err != nil {
		return err
	}

	if config == nil {
		app.Config = &Config{}
	} else {
		app.Config = config
	}

	return nil
}

func defaultConfigLoader(file string) (*Config, error) {
	if file == "" {
		return &Config{}, nil
	}
	return LoadConfigFile(file)
}

func getContextConfigFile(c *cli.Context) string {
	file := c.String("conf")
	if file == "" {
		file = findConfigFile(c.App.Name)
	}
	return file
}

func findConfigFile(appName string) string {
	filename := appName + ".toml"

	// current directory
	if dir, err := os.Getwd(); err == nil {
		if file := findConfigFileInDir(dir, "config", filename); file != "" {
			return file
		}
	}

	// executable directory
	if dir, err := osext.ExecutableFolder(); err == nil {
		if file := findConfigFileInDir(dir, "config", filename); file != "" {
			return file
		}
	}

	// home directory
	// TODO: support XDG_CONFIG_HOME and XDG_CONFIG_DIRS
	if dir, err := homedir.Dir(); err == nil {
		subdir := filepath.Join(".config", appName)
		if file := findConfigFileInDir(dir, subdir, filename); file != "" {
			return file
		}
	}

	return ""
}

func findConfigFileInDir(dir, subdir, filename string) string {
	file := filepath.Join(dir, filename)
	if fileExists(file) {
		return file
	}

	// under sub-directory
	if subdir != "" {
		file = filepath.Join(dir, subdir, filename)
		if fileExists(file) {
			return file
		}
	}

	return ""
}
