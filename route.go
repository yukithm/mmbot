package mmbot

import "net/http"

// RouteHandlerFunc is route action function.
type RouteHandlerFunc func(*Robot, http.ResponseWriter, *http.Request) error

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
