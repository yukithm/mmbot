package mmbot

import (
	"log"
	"mmbot/mmhook"
	"net/http"

	"github.com/gorilla/mux"
)

type Robot struct {
	*mmhook.Client
	Config   *Config
	Handlers []Handler
	quit     chan bool
}

func NewRobot(config *Config) *Robot {
	bot := &Robot{
		Client: mmhook.NewClient(config.Config),
		Config: config,
		quit:   make(chan bool),
	}

	return bot
}

func (r *Robot) Run() {
	if !r.Config.DisableServer {
		r.StartServer()
		r.receive()
	} else {
		<-r.quit
	}
}

func (r *Robot) Stop() {
	r.quit <- true
}

func (r *Robot) receive() {
	for {
		select {
		case <-r.quit:
			return
		case msg := <-r.In:
			r.handle(&msg)
		}
	}
}

func (r *Robot) handle(inMsg *mmhook.InMessage) {
	msg := &InMessage{
		InMessage: inMsg,
		Robot:     r,
	}

	for _, handler := range r.Handlers {
		if handler.CanHandle(msg) {
			err := handler.Handle(msg)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func (r *Robot) startServer() {
	mux := mux.NewRouter()
	r.mountClient(mux)

	// TODO: mount http handler

	log.Printf("Listening on %s\n", r.Address())
	if err := http.ListenAndServe(r.Address(), mux); err != nil {
		log.Fatal(err)
	}
}

func (r *Robot) mountClient(mux *mux.Router) {
	path := r.Client.IncomingPath
	if path == "" {
		path = "/"
	}
	mux.Handle(path, r.Client)
}
