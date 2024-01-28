package api

import (
	"net/http"
)

// DefaultHandler handles the default route.
func DefaultHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
