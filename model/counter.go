package model

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jpcercal/counter/config"
)

// Counter is the interface for the counter.
var Counter SafeCounter

// TimeWindowInSeconds is the time window in which requests are counted.
var TimeWindowInSeconds = config.ServerTimeWindowInSecondsToCountRequests

// filename is the name of the file where the counter state is stored.
const filename = config.ServerFilenameToStoreCounterState

// safeCounter is safe to use concurrently.
type SafeCounter struct {
	mu          sync.Mutex
	lastRequest []time.Time
}

// Init initializes the counter.
func (s *SafeCounter) Init() {
	s.LoadState()
}

// LoadState loads the counter state from disk.
func (s *SafeCounter) LoadState() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// load serialized state from disk.
	log.Println("loading state from disk...")
	log.Printf("the counter time window is %d seconds...\n", TimeWindowInSeconds)

	// TODO: check if some other format is better to store the state.
	//
	//nolint:godox

	// read the file from disk.
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("error opening state file: %s\n", err.Error())
	}

	// decode the json data.
	if err := json.Unmarshal(jsonData, &s.lastRequest); err != nil {
		log.Printf("error decoding state: %s\n", err.Error())
	}
}

// ResetState resets the counter state.
func (s *SafeCounter) ResetState() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("reseting state...")

	s.lastRequest = []time.Time{}
}

// SaveState saves the current state of the counter to disk.
func (s *SafeCounter) SaveState() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("saving state to disk on a file named %s...\n", filename)

	jsonData, err := json.Marshal(s.lastRequest)
	if err != nil {
		log.Printf("error encoding state: %s\n", err.Error())
	}

	const fileMode os.FileMode = 0600

	err = os.WriteFile(filename, jsonData, fileMode)
	if err != nil {
		log.Printf("error creating state file: %s\n", err.Error())
	}
}

// increment increments the counter.
func (s *SafeCounter) Increment() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastRequest = append(s.lastRequest, time.Now())
}

// value returns the current value of the counter.
func (s *SafeCounter) Value() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	// clean up requests older than the time window.
	count := 0

	// count the requests older than the time window.
	for i := 0; i < len(s.lastRequest); i++ {
		if time.Since(s.lastRequest[i]).Seconds() > float64(TimeWindowInSeconds) {
			count++
		}
	}

	// remove the requests older than the time window.
	s.lastRequest = s.lastRequest[count:]

	return len(s.lastRequest)
}
