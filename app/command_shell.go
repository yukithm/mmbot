package app

import (
	"mmbot"
	"mmbot/shell"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
)

func (app *App) newShellCommand() cli.Command {
	return cli.Command{
		Name:        "shell",
		Usage:       "run interactive shell",
		Description: "run interactive shell",
		Flags: []cli.Flag{
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
		Action: app.shellCommand,
	}
}

func (app *App) shellCommand(c *cli.Context) {
	app.updateConfigByFlags(c)
	app.Config.ValidateAndExitOnError()

	client := shell.NewClient(app.Config.AdapterConfig(), app.Config.Logger)
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