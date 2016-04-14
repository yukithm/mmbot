package mmbot

import "net/http"

type RouteHandlerFunc func(*Robot, http.ResponseWriter, *http.Request) error

type Route struct {
	Methods []string
	Pattern string
	Action  RouteHandlerFunc
}
