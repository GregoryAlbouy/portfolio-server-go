package utl

import (
	"net/http"
)

// RouteHandler type
type RouteHandler func(http.ResponseWriter, *http.Request)

// RouteMap type
type RouteMap map[string]RouteHandler
