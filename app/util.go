package app

import (
	"os"
	"path/filepath"

	"github.com/VividCortex/godaemon"
	"github.com/codegangsta/cli"
)

func (app *App) updateConfigByFlags(c *cli.Context) {
	if c.IsSet("log") {
		app.Config.Common.Log = c.String("log")
	}
	if c.IsSet("daemonize") {
		app.Config.Common.daemonize = c.Bool("daemonize")
	}
	if c.IsSet("pidfile") {
		app.Config.Common.PIDFile = c.String("pidfile")
	}
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

func (app *App) newLogger() (*Logger, error) {
	c := app.Config.Common
	if c.daemonize && (c.Log == "" || c.Log == "-") {
		return NewNullLogger()
	}
	return NewLogger(c.Log)
}

func absPath(file string) (string, error) {
	if file == "" {
		return "", nil
	}

	if !filepath.IsAbs(file) {
		exec, err := godaemon.GetExecutablePath()
		if err != nil {
			return "", err
		}
		file = filepath.Join(filepath.Dir(exec), file)
	}

	return file, nil
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}
