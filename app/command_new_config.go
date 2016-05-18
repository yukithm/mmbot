package app

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/codegangsta/cli"
)

func (app *App) newNewConfigCommand() cli.Command {
	return cli.Command{
		Name:        "new-config",
		Usage:       "Create new config file",
		Description: "Create new config file",
		Action:      app.newConfigCommand,
	}
}

func (app *App) newConfigCommand(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	filename := filepath.Join(cwd, c.App.Name+".toml")

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return cli.NewExitError(fmt.Sprintf("%s already exists.", filename), 1)
		}
		return cli.NewExitError(err.Error(), 1)
	}
	defer file.Close()

	vars := struct {
		Name string
	}{
		Name: c.App.Name,
	}
	tmpl := template.Must(template.New("config").Parse(NewConfigTemplate))
	if err := tmpl.Execute(file, vars); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Printf("Written new config file to %s\n", filename)

	return nil
}

// NewConfigTemplate is a template for configuration file.
// You can change it content by adding your extra configuration entries.
// {{.Name}} is replaced with the application name.
var NewConfigTemplate = `# {{.Name}} configuration file
#
# TOML format
# See: https://github.com/toml-lang/toml

[common]
# Log file path (empty: STDERR, "-": STDOUT)
# log = "./{{.Name}}.log"

# PID file path (empty: not create)
# pidfile = "/var/run/{{.Name}}.pid"

[mattermost]
# Webhook URL for posting messages (REQUIRED)
# (send to Mattermost; Incoming Webhooks on Mattermost side)
outgoing_url = "http://localhost/incoming_webhook_url"

# Webhook path on the bot for receiving messages (default: "/")
# (receive from Mattermost; Outgoing Webhooks on Mattermost side)
# NOTE: You need to enable HTTP server (server.enable = true)
incoming_path = "/{{.Name}}_incoming"

# Tokens from Mattermost outgoing webhooks (default: [])
# If omitted, all requests are accepted
tokens = [
    "incomign_webhook_token"
]

# Username of the bot account (preceded by '@') (REQUIRED)
username = "{{.Name}}"

# Overridding of username for Mattermost webhook (default: "")
# override_username = "{{.Name}}"

# Overridding of icon URL for Mattermost webhook (default: "")
# icon_url = "http://localhost/{{.Name}}.png"

# Disable certificate checking (default: false)
# insecure_skip_verify = true

[server]
# Enable HTTP server for webhook and handlers (default: false)
enable = true

# Bind address for the bot HTTP server (default: ""; all interfaces)
bind_address = ""

# Bind port for the bot HTTP server (default: 8080)
port = 8080
`
