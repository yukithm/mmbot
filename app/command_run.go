package app

import (
	"mmbot"
	"mmbot/mmhook"
	"os"
	"os/signal"
	"syscall"

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
			cli.StringSliceFlag{
				Name:  "tokens",
				Usage: "tokens from Mattermost outgoing webhooks",
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

	client := mmhook.NewClient(app.Config.AdapterConfig(), app.Config.Logger)
	robot := mmbot.NewRobot(app.Config.RobotConfig(), client)
	robot.Handlers = app.Handlers
	robot.Routes = app.Routes
	robot.Jobs = app.Jobs

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		s := <-sigCh
		app.Config.Logger.Printf("%q received", s)
		robot.Stop()
	}()

	robot.Run()
	app.Config.Logger.Println("Stop robot")
}
