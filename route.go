package mmbot

import (
	"net/http"

	"github.com/fukata/golang-stats-api-handler"
)

// RouteHandlerFunc is route action function.
type RouteHandlerFunc func(*Robot, http.ResponseWriter, *http.Request)

// Route is a HTTP route.
type Route struct {
	// HTTP method (e.g. "GET", "POST").
	// All methods are allowed if empty.
	Methods []string

	// Route pattern (e.g. "/articles/{category}/{id:[0-9]+}")
	// Pattern can have variables that are defined by "{name}" or "{name:regexp}" format.
	// Variables can be retrieved calling Robot.RouteVars().
	Pattern string

	// Route action.
	Action RouteHandlerFunc
}

// NewPingRoute returns the route "ping".
func NewPingRoute(pattern string) Route {
	return Route{
		Methods: []string{"GET"},
		Pattern: pattern,
		Action: func(bot *Robot, w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		},
	}
}

// NewStatsRoute returns the route for statistics of the process.
func NewStatsRoute(pattern string) Route {
	return Route{
		Methods: []string{"GET"},
		Pattern: pattern,
		Action: func(bot *Robot, w http.ResponseWriter, r *http.Request) {
			stats_api.Handler(w, r)
		},
	}
}
