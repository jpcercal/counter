package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/jpcercal/counter/api"
	"github.com/jpcercal/counter/model"
	testing_helper "github.com/jpcercal/counter/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCounterHandler tests the counter handler.
//
// TODO: the usage of testcontainers is a good idea to test the counter,
// and to run the tests in parallel. I would probably have used it to make it
// the closest to a real scenario.
//
//nolint:godox
func TestCounterHandler(t *testing.T) {
	// makeRequest is a helper function to make a request to the server
	makeRequest := func(urlToRequest string) int {
		request, err := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			urlToRequest,
			nil,
		)
		require.NoError(t, err)

		// Send the request
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)

		// Close the response body when test finishes
		defer response.Body.Close()

		// Check the response status code
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// Check the response content type
		assert.Equal(t, "application/json", response.Header.Get("Content-Type"))

		// Read the response body
		body, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		// Check the response body
		var counterResponse api.CounterResponse

		// Check the json can be unmarshalled
		require.NoError(t, json.Unmarshal(body, &counterResponse))

		return counterResponse.Counter
	}

	t.Run("test if the counter works with multiple concurrent requests", func(t *testing.T) {
		// Reset the counter state
		model.Counter.ResetState()

		// Set the time window to 30 seconds
		model.TimeWindowInSeconds = 30

		// Get a free TCP port
		port, err := testing_helper.GetFreeTCPPort()
		require.NoError(t, err)

		// Create a new server
		server, err := api.NewServer(port)
		require.NoError(t, err)

		// Close the server when test finishes
		defer (func() {
			t.Log("Shutting down the server gracefully...")

			err := server.Shutdown(context.Background())

			require.NoError(t, err)
		})()

		// Start the server in a goroutine
		go server.Start()

		// Check if the server is running
		url := fmt.Sprintf("http://localhost:%d/%s", port, api.CounterRoutePath)

		// Check the amount of requests made with the counter
		assert.Equal(t, 1, makeRequest(url))

		// numer of requests to be fired on each goroutine
		numberOfRequestsToBeFired := 1000

		t.Logf("Number of CPUs: %d", runtime.NumCPU())
		t.Logf("Number of requests to be fired on each goroutine: %d", numberOfRequestsToBeFired)

		// use a waitgroup
		var wg sync.WaitGroup

		// add the number of goroutines to the waitgroup
		wg.Add(numberOfRequestsToBeFired * runtime.NumCPU())

		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for j := 0; j < numberOfRequestsToBeFired; j++ {
					makeRequest(url)
					wg.Done()
				}
			}()
		}

		// wait for all goroutines to finish
		wg.Wait()

		// Check the amount of requests made with the counter
		assert.LessOrEqual(t, (numberOfRequestsToBeFired*runtime.NumCPU())+1, model.Counter.Value())
	})

	t.Run("test if the counter works after passing the time window", func(t *testing.T) {
		// Reset the counter state
		model.Counter.ResetState()

		// Set the time window to 2 seconds
		model.TimeWindowInSeconds = 2

		// Get a free TCP port
		port, err := testing_helper.GetFreeTCPPort()
		require.NoError(t, err)

		// Create a new server
		server, err := api.NewServer(port)
		require.NoError(t, err)

		// Close the server when test finishes
		defer (func() {
			t.Log("Shutting down the server gracefully...")

			err := server.Shutdown(context.Background())

			require.NoError(t, err)
		})()

		// Start the server in a goroutine
		go server.Start()

		// Check if the server is running
		url := fmt.Sprintf("http://localhost:%d/%s", port, api.CounterRoutePath)

		// Fire 10 requests and check the counter
		// value after each request
		assert.Equal(t, 1, makeRequest(url))
		assert.Equal(t, 2, makeRequest(url))
		assert.Equal(t, 3, makeRequest(url))
		assert.Equal(t, 4, makeRequest(url))
		assert.Equal(t, 5, makeRequest(url))
		assert.Equal(t, 6, makeRequest(url))
		assert.Equal(t, 7, makeRequest(url))
		assert.Equal(t, 8, makeRequest(url))
		assert.Equal(t, 9, makeRequest(url))
		assert.Equal(t, 10, makeRequest(url))

		// Wait for 3 seconds to pass the time window
		time.Sleep(3 * time.Second)
		assert.Equal(t, 1, makeRequest(url))

		// Wait for 3 seconds to pass the time window
		time.Sleep(3 * time.Second)
		makeRequest(url)
		makeRequest(url)
		assert.Equal(t, 3, makeRequest(url))
	})

	t.Run("test if the counter works after passing the time window and the server is restarted", func(t *testing.T) {
		// Reset the counter state
		model.Counter.ResetState()

		// Set the time window to 10 seconds
		model.TimeWindowInSeconds = 10

		// Get a free TCP port
		port, err := testing_helper.GetFreeTCPPort()
		require.NoError(t, err)

		// Create a new server
		server, err := api.NewServer(port)
		require.NoError(t, err)

		// Check if the server is running
		url := fmt.Sprintf("http://localhost:%d/%s", port, api.CounterRoutePath)

		// Close the server when test finishes
		defer (func() {
			t.Log("Shutting down the server gracefully for the first time...")

			// Shutdown the server gracefully
			err := server.Shutdown(context.Background())
			require.NoError(t, err)

			// Get a free TCP port
			port, err := testing_helper.GetFreeTCPPort()
			require.NoError(t, err)

			// Create a new server
			server, err := api.NewServer(port)
			require.NoError(t, err)

			// Check if the server is running
			url := fmt.Sprintf("http://localhost:%d/%s", port, api.CounterRoutePath)

			// Start the server in a goroutine
			go server.Start()

			// Close the server when test finishes
			defer (func() {
				t.Log("Shutting down the server gracefully for the second time...")
				err := server.Shutdown(context.Background())
				require.NoError(t, err)
			})()

			// Make sure that the counter is not reseted
			assert.Equal(t, 10, model.Counter.Value())

			// Wait for 2 seconds, it will not pass the time window
			time.Sleep(2 * time.Second)
			makeRequest(url)
			makeRequest(url)

			// It must be 13 because the counter was not reseted
			// after the server was restarted
			assert.Equal(t, 13, makeRequest(url))
		})()

		// Start the server in a goroutine
		go server.Start()

		// Fire 10 requests and check the counter
		// value after each request
		assert.Equal(t, 1, makeRequest(url))
		assert.Equal(t, 2, makeRequest(url))
		assert.Equal(t, 3, makeRequest(url))
		assert.Equal(t, 4, makeRequest(url))
		assert.Equal(t, 5, makeRequest(url))
		assert.Equal(t, 6, makeRequest(url))
		assert.Equal(t, 7, makeRequest(url))
		assert.Equal(t, 8, makeRequest(url))
		assert.Equal(t, 9, makeRequest(url))
		assert.Equal(t, 10, makeRequest(url))
	})
}
