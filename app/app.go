// Package app provides a base of the bot application.
// App handles command line arguments and manages mmbot.Robot.
package app

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/yukithm/mmbot"
)

// App is a bot application.
type App struct {
	*cli.App
	Config   *Config
	Handlers []mmbot.Handler
	Routes   []mmbot.Route
	Jobs     []mmbot.Job
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
		app.newRunCommand(),
		app.newShellCommand(),
	}

	app.App.Before = func(c *cli.Context) error {
		config, err := loadConfig(c)
		if err != nil {
			log.Fatal(err)
		}
		app.Config = config
		return nil
	}

	return app
}
