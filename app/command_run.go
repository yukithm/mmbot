package app

import (
	"log"
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
			cli.StringFlag{
				Name:  "log",
				Usage: "log file",
			},
		},
		Action: app.runCommand,
	}
}

func (app *App) runCommand(c *cli.Context) {
	app.updateConfigByFlags(c)
	app.Config.ValidateAndExitOnError()

	logger, err := app.newLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	client := mmhook.NewClient(app.Config.AdapterConfig(), logger.Logger)
	robot := mmbot.NewRobot(app.Config.RobotConfig(), client, logger.Logger)
	robot.Handlers = app.Handlers
	robot.Routes = app.Routes
	robot.Jobs = app.Jobs

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	quit := make(chan bool)

	go func() {
		s := <-sigCh
		logger.Printf("%q received", s)
		close(quit)
	}()

	errCh := robot.Start()

	select {
	case <-quit:
		robot.Stop()
		logger.Println("Stop robot")
	case <-errCh:
		logger.Println("Abort robot")
	}
}
