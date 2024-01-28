package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jpcercal/counter/model"
)

// CounterRoutePath is the route path for the counter.
const CounterRoutePath = "/counter"

// CounterResponse is the http response that represents the counter.
type CounterResponse struct {
	Counter int `json:"counter"`
}

// CounterHandler handles the counter route.
func CounterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	// Increment the counter (thread-safe)
	model.Counter.Increment()

	// Return the counter value
	resp := CounterResponse{
		Counter: model.Counter.Value(), // read the counter (thread-safe)
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error encoding response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	// Set the response headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(jsonData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
