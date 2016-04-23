package mmbot

import (
	"fmt"
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
	Config     *Config
	Client     adapter.Adapter
	Handlers   []Handler
	Routes     []Route
	Jobs       []Job
	scheduler  *cron.Cron
	logger     *log.Logger
	workerJobs chan workerJob
	aborted    bool
	quit       chan struct{}
	errCh      chan error
}

type workerJob struct {
	handler Handler
	message *message.InMessage
}

const (
	numJobWorkers = 4
	numJobBuffers = 20
)

func NewRobot(config *Config, client adapter.Adapter, logger *log.Logger) *Robot {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	bot := &Robot{
		Config: config,
		Client: client,
		logger: logger,
	}

	return bot
}

func (r *Robot) Start() chan error {
	r.errCh = make(chan error, 1)
	go r.run()
	return r.errCh
}

func (r *Robot) run() {
	r.aborted = false
	r.quit = make(chan struct{}, 1)

	r.workerJobs = make(chan workerJob, numJobBuffers)
	for i := 1; i <= numJobWorkers; i++ {
		go r.worker(i, r.workerJobs)
	}

	r.runLoop()

	if !r.aborted {
		r.Client.Stop()
		r.logger.Println("Stop adapter")
	}

	if r.scheduler != nil {
		r.scheduler.Stop()
		r.logger.Println("Stop job scheduler")
	}

	close(r.workerJobs)
	close(r.quit)
	close(r.errCh)
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
				r.aborted = true
				r.logger.Print(e)
			}
			return
		case msg, ok := <-receiver:
			if !ok {
				r.aborted = true
				return
			}
			r.handle(&msg)
		}
	}
}

func (r *Robot) Stop() {
	if !r.aborted {
		r.quit <- struct{}{}
		<-r.quit
	}
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
		r.workerJobs <- workerJob{
			handler: handler,
			message: msg,
		}
	}
}

func (r *Robot) worker(id int, jobs <-chan workerJob) {
	for job := range jobs {
		r.callHandler(job.handler, job.message)
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
		r.errCh <- err
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
			r.errCh <- fmt.Errorf("Invalid route: %v", route)
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
