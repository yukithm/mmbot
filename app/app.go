package app

import (
	"log"

	"github.com/codegangsta/cli"
)

// App is a bot application.
type App struct {
	*cli.App
	Config *Config
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
