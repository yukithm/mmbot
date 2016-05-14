package app

import (
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func (app *App) updateConfigByFlags(c *cli.Context) {
	if c.IsSet("log") {
		app.Config.Common.Log = c.String("log")
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

// Logger is a logger that has *os.File.
type Logger struct {
	*log.Logger
	file *os.File
}

// Close close the log file when it is not nil.
func (l *Logger) Close() error {
	if l.file != nil {
		err := l.file.Close()
		if err != nil {
			return err
		}
		l.file = nil
	}
	return nil
}

func (app *App) newLogger() (*Logger, error) {
	var file *os.File
	var w io.Writer

	if app.Config.Common.Log == "" {
		w = os.Stderr
	} else if app.Config.Common.Log == "-" {
		w = os.Stdout
	} else {
		file, err := os.OpenFile(app.Config.Common.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		w = file
	}

	logger := &Logger{
		Logger: log.New(w, "", log.LstdFlags),
		file:   file,
	}
	return logger, nil
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}
