package mmbot

import (
	"io/ioutil"
	"log"
	"mmbot/adapter"
	"mmbot/message"
	"net/http"
	"runtime"
	"time"

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

func NewRobot(config *Config, client adapter.Adapter) *Robot {
	if config.Logger == nil {
		config.Logger = log.New(ioutil.Discard, "", 0)
	}
	bot := &Robot{
		Config: config,
		Client: client,
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

	receiver, errCh := r.Client.Start()

	r.startScheduler()

	for {
		select {
		case <-r.quit:
			return
		case e, ok := <-errCh:
			if ok {
				r.logger.Print(e)
			}
			return
		case msg, ok := <-receiver:
			if !ok {
				return
			}
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
		r.callHandler(handler, msg)
	}
}

func (r *Robot) callHandler(handler Handler, msg *message.InMessage) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			r.logger.Printf("mmbot: panic handler: %v\n%s", err, buf)
		}
	}()

	if handler.CanHandle(msg) {
		err := handler.Handle(msg)
		if err != nil {
			r.logger.Print(err)
		}
	}
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
	server := &http.Server{
		Addr:        r.Config.Address(),
		Handler:     mux,
		ReadTimeout: 30 * time.Second,
		ErrorLog:    r.logger,
	}
	if err := server.ListenAndServe(); err != nil {
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
