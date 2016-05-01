package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/yukithm/mmbot"
	"github.com/yukithm/mmbot/mmhook"
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
		Before: func(c *cli.Context) error {
			if err := app.LoadConfig(c); err != nil {
				log.Fatal(err)
			}
			return nil
		},
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

	if app.InitRobot != nil {
		app.InitRobot(robot)
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	quit := make(chan struct{})

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
	case err, ok := <-errCh:
		if ok && err != nil {
			logger.Print(err)
			logger.Println("Abort robot")
		} else {
			logger.Println("Stop robot")
		}
	}
}
