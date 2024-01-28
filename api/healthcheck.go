package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HealthCheckRoutePath is the route path for the health check.
const HealthCheckRoutePath = "/healthcheck"

// HealthCheckResponse is the http response that represents the health check.
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// HealthCheckHandler handles the health check route.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	resp := HealthCheckResponse{
		Status:    "OK",
		Timestamp: time.Now().UnixMilli(),
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
