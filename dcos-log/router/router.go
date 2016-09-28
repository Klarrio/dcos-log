package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter returns a new instance of gorilla.mux with loaded routes.
func NewRouter(routes []Route) (*mux.Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		if route.URL == "" {
			return router, fmt.Errorf("URL cannot be empty for the route: %v", route)
		}

		if route.Handler == nil {
			return router, fmt.Errorf("Handler function cannot be empty for the route: %v", route)
		}

		if len(route.Headers) != 2 {
			return router, fmt.Errorf("Headers length must be 2, route: %v", route)
		}

		handler := http.HandlerFunc(route.Handler)
		router.Handle(route.URL, wrapHandler(handler, route)).Methods("GET").HeadersRegexp(route.Headers...)
	}
	return router, nil
}

// Route is a structure defines a set of parameters passed to gorilla.mux router via .Handle method
type Route struct {
	// URL is endpoint url
	URL string

	// Handler is a handler function
	Handler func(http.ResponseWriter, *http.Request)

	// Headers is pair of headers passed to HeadersRegexp
	Headers []string
}

func wrapHandler(handler http.Handler, route Route) http.Handler {
	// apply appropriate middleware
	handler = validateQueryStringMiddleware(authMiddleware(handler))
	return handler
}
