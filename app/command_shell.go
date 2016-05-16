package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/yukithm/mmbot"
	"github.com/yukithm/mmbot/shell"
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
			cli.StringFlag{
				Name:  "log",
				Usage: "log file",
			},
		},
		Action: app.shellCommand,
		Before: func(c *cli.Context) error {
			if err := app.LoadConfig(c); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
	}
}

func (app *App) shellCommand(c *cli.Context) error {
	app.updateConfigByFlags(c)
	app.Config.ValidateAndExitOnError()

	logger, err := app.newLogger()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer logger.Close()

	client := shell.NewClient(app.Config.AdapterConfig(), logger.Logger)
	robot := mmbot.NewRobot(app.Config.RobotConfig(), client, logger.Logger)

	if app.InitRobot != nil {
		if err := app.InitRobot(robot); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
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

	return nil
}
