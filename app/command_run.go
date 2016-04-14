package app

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func (app *App) newRunCommand() cli.Command {
	return cli.Command{
		Name:        "run",
		Usage:       "start bot",
		Description: "start bot",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "outgoing-url",
				Usage: "webhook URL for Mattermost (Incoming Webhooks on Mattermost side)",
			},
			cli.StringFlag{
				Name:  "incoming-path",
				Value: "/",
				Usage: "webhook path from Mattermost (Outgoing Webhooks on Mattermost side)",
			},
			cli.StringFlag{
				Name:  "token",
				Usage: "toke from Mattermost outgoing webhooks",
			},
			cli.StringFlag{
				Name:  "username",
				Usage: "username of the bot account",
			},
			cli.StringFlag{
				Name:  "override-username",
				Usage: "overriding of username",
			},
			cli.StringFlag{
				Name:  "icon-url",
				Usage: "overriding of icon URL",
			},
			cli.BoolFlag{
				Name:  "insecure-skip-verify",
				Usage: "disable certificate checking",
			},
			cli.BoolFlag{
				Name:  "disable-server",
				Usage: "disable the bot HTTP server",
			},
			cli.StringFlag{
				Name:  "bind-address",
				Usage: "bind address for the bot HTTP server",
			},
			cli.IntFlag{
				Name:  "port",
				Value: 8080,
				Usage: "bind port for the bot HTTP server",
			},
		},
		Action: app.runCommand,
	}
}

func (app *App) runCommand(c *cli.Context) {
	app.updateConfigByFlags(c)
	app.Config.ValidateAndExitOnError()
	fmt.Printf("%#v\n", app.Config)
}

func (app *App) updateConfigByFlags(c *cli.Context) {
	if c.IsSet("outgoing-url") {
		app.Config.Mattermost.OutgoingURL = c.String("outgoing-url")
	}
	if c.IsSet("incoming-path") {
		app.Config.Mattermost.IncomingPath = c.String("incoming-path")
	}
	if c.IsSet("token") {
		app.Config.Mattermost.Token = c.String("token")
	}
	if c.IsSet("username") {
		app.Config.Mattermost.UserName = c.String("username")
	}
	if c.IsSet("override-username") {
		app.Config.Mattermost.OverrideUserName = c.String("override-username")
	}
	if c.IsSet("icon-url") {
		app.Config.Mattermost.IconURL = c.String("icon-url")
	}
	if c.IsSet("insecure-skip-verify") {
		app.Config.Mattermost.InsecureSkipVerify = c.Bool("insecure-skip-verify")
	}
	if c.IsSet("disable-server") {
		app.Config.Server.Enable = !c.Bool("disable-server")
	}
	if c.IsSet("bind-address") {
		app.Config.Server.BindAddress = c.String("bind-address")
	}
	if c.IsSet("port") {
		app.Config.Server.Port = c.Int("port")
	}
}
