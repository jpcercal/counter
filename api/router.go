package api

import (
	"net/http"
)

// New returns an http.Handler configured with application routes.
func NewHTTPHandler() (http.Handler, error) {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", DefaultHandler)
	mux.HandleFunc(CounterRoutePath, CounterHandler)
	mux.HandleFunc(HealthCheckRoutePath, HealthCheckHandler)

	return mux, nil
}
