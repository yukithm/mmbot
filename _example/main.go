package main

import (
	"fmt"
	"github.com/yukithm/mmbot"
	"github.com/yukithm/mmbot/app"
	"github.com/yukithm/mmbot/message"
	"net/http"
	"regexp"
	"time"
)

const (
	// ApplicationName is the name of this application.
	ApplicationName = "mmbot"

	// Version is the version number of this application.
	Version = "0.1.0"
)

func main() {
	app := app.NewApp()
	app.Name = ApplicationName
	app.Version = Version
	app.Usage = "a bot for Mattermost"

	app.Handlers = []mmbot.Handler{
		mmbot.PatternHandler{
			Pattern: regexp.MustCompile(`\Ahello`),
			Action: func(msg *message.InMessage) error {
				// raw := msg.RawMessage.(*mmhook.InMessage)
				// fmt.Printf("msg=%#v, raw=%#v\n", msg, raw)
				return msg.Reply("Hello, " + msg.UserName)
			},
		},
		mmbot.PatternHandler{
			Pattern: regexp.MustCompile(`\Aこんにち[はわ]`),
			Action: func(msg *message.InMessage) error {
				fmt.Printf("msg=%#v\n", msg)
				// msg.Sender.Send(&message.OutMessage{
				// 	ChannelName: "town-square",
				// 	Text:        msg.UserName + "さんに挨拶しました",
				// })
				// msg.Sender.Send(&message.OutMessage{
				// 	ChannelName: "@" + msg.UserName,
				// 	Text:        msg.UserName + "さん、よろしくお願いします",
				// })
				return msg.Reply(fmt.Sprintf("@%sさん、こんにちは！", msg.UserName))
			},
		},
	}

	app.Routes = []mmbot.Route{
		mmbot.NewPingRoute("/ping"),
		mmbot.NewStatsRoute("/stats"),
		mmbot.Route{
			Methods: []string{"GET"},
			Pattern: "/hello",
			Action: func(bot *mmbot.Robot, w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello!"))
			},
		},
	}

	app.Jobs = []mmbot.Job{
		mmbot.Job{
			Schedule: "0 * * * * *",
			Action: func(bot *mmbot.Robot) {
				fmt.Printf("job: %s", time.Now())
				bot.Send(&message.OutMessage{
					Text: fmt.Sprintf("job: %s", time.Now()),
				})
			},
		},
	}

	app.RunAndExitOnError()
}
