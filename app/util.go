package app

import (
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func (app *App) updateConfigByFlags(c *cli.Context) {
	if c.IsSet("outgoing-url") {
		app.Config.Mattermost.OutgoingURL = c.String("outgoing-url")
	}
	if c.IsSet("incoming-path") {
		app.Config.Mattermost.IncomingPath = c.String("incoming-path")
	}
	if c.IsSet("tokens") {
		app.Config.Mattermost.Tokens = c.StringSlice("token")
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

func (app *App) initLogger(logfile string) error {
	var w io.Writer
	if logfile == "" {
		w = os.Stderr
	} else if logfile == "-" {
		w = os.Stdout
	} else {
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		app.Config.logfile = f
		w = f
	}
	app.Config.Logger = log.New(w, "", log.LstdFlags)
	return nil
}

func (app *App) closeLogger() error {
	if app.Config.logfile != nil {
		err := app.Config.logfile.Close()
		if err != nil {
			return err
		}
		app.Config.logfile = nil
	}
	return nil
}
