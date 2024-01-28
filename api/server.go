package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jpcercal/counter/model"
)

// Server represents the http server.
type Server struct {
	*http.Server
}

// Server timeouts, to avoid Slowloris attacks.
// G112: Potential Slowloris Attack because ReadHeaderTimeout is not configured in the http.Server.
const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
)

// NewServer returns a new instance of Server.
func NewServer(port int) (*Server, error) {
	log.Println("configuring server...")

	// Create the http handler
	api, err := NewHTTPHandler()
	if err != nil {
		return nil, fmt.Errorf("error creating http handler: %w", err)
	}

	// Create the http server
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      api,
	}

	return &Server{&srv}, nil
}

// Start starts the server.
func (srv *Server) Start() {
	log.Println("starting server...")

	// Initialize the counter.
	// TODO: move this to a better place Using a router like gin or chi, we could
	// use a middleware to initialize the counter to be done.
	//
	//nolint:godox
	model.Counter.Init()
	defer model.Counter.SaveState()

	// Start the server
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	// Gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)

	// Shutdown the server gracefully
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	log.Println("Server gracefully stopped")
}
