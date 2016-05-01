// Package app provides a base of the bot application.
// App handles command line arguments and manages mmbot.Robot.
package app

import (
	"github.com/codegangsta/cli"
	"github.com/yukithm/mmbot"
)

// App is a bot application.
type App struct {
	*cli.App
	Config    *Config
	InitRobot func(*mmbot.Robot)
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
	config, err := loadConfig(c)
	if err != nil {
		return err
	}
	app.Config = config
	return nil
}
