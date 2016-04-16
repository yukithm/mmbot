package mmbot

import (
	"io/ioutil"
	"log"
	"mmbot/adapter"
	"mmbot/message"
	"mmbot/mmhook"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

type Robot struct {
	Config    *Config
	Client    adapter.Adapter
	Handlers  []Handler
	Routes    []Route
	Jobs      []Job
	scheduler *cron.Cron
	logger    *log.Logger
	quit      chan bool
}

func NewRobot(config *Config) *Robot {
	if config.Logger == nil {
		config.Logger = log.New(ioutil.Discard, "", 0)
	}
	bot := &Robot{
		Config: config,
		Client: mmhook.NewClient(config.AdapterConfig, config.Logger),
		logger: config.Logger,
		quit:   make(chan bool),
	}

	return bot
}

func (r *Robot) Run() {
	r.runLoop()

	r.Client.Stop()
	r.logger.Println("Stop adapter")

	if r.scheduler != nil {
		r.scheduler.Stop()
		r.logger.Println("Stop job scheduler")
	}
}

func (r *Robot) runLoop() {
	if !r.Config.DisableServer {
		go r.startServer()
	}
	go r.startClient()

	receiver := r.Client.Receiver()
	for {
		select {
		case <-r.quit:
			return
		case msg := <-receiver:
			r.handle(&msg)
		}
	}
}

func (r *Robot) Stop() {
	r.quit <- true
}

func (r *Robot) Send(msg *message.OutMessage) error {
	return r.Client.Send(msg)
}

func (r *Robot) SenderName() string {
	return r.Config.UserName
}

func (r *Robot) handle(msg *message.InMessage) {
	msg.Sender = r

	for _, handler := range r.Handlers {
		if handler.CanHandle(msg) {
			err := handler.Handle(msg)
			if err != nil {
				r.logger.Print(err)
			}
		}
	}
}

func (r *Robot) startClient() {
	err := r.Client.Run()
	if err != nil {
		r.logger.Print(err)
	}
	r.Stop()
}

func (r *Robot) startScheduler() {
	if r.Jobs == nil || len(r.Jobs) == 0 {
		return
	}

	r.scheduler = cron.New()
	for _, job := range r.Jobs {
		r.scheduler.AddFunc(job.Schedule, func() {
			job.Action(r)
		})
	}
	r.scheduler.Start()
	r.logger.Println("Start job scheduler")
}

func (r *Robot) startServer() {
	mux := mux.NewRouter()
	r.mountRoutes(mux)
	r.mountClient(mux)

	r.logger.Printf("Listening on %s\n", r.Config.Address())
	if err := http.ListenAndServe(r.Config.Address(), mux); err != nil {
		r.logger.Fatal(err)
	}
	r.Stop()
}

func (r *Robot) mountClient(mux *mux.Router) {
	hook := r.Client.IncomingWebHook()
	if hook == nil {
		return
	}

	path := hook.Path
	if path == "" {
		path = "/"
	}
	mux.Handle(path, hook.Handler)
}

func (r *Robot) mountRoutes(mux *mux.Router) {
	for _, route := range r.Routes {
		if route.Pattern == "" || route.Action == nil {
			log.Fatalf("Invalid route: %v", route)
		}

		wrapped := r.wrapRouteAction(route)
		mr := mux.HandleFunc(route.Pattern, wrapped)
		if route.Methods != nil && len(route.Methods) > 0 {
			mr.Methods(route.Methods...)
		}
	}
}

func (r *Robot) wrapRouteAction(route Route) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		route.Action(r, w, req)
	}
}

func (r *Robot) RouteVars(req *http.Request) map[string]string {
	return mux.Vars(req)
}
