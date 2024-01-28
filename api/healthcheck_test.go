package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jpcercal/counter/api"
	testing_helper "github.com/jpcercal/counter/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealthCheckHandler tests the health check handler.
func TestHealthCheckHandler(t *testing.T) {
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
	url := fmt.Sprintf("http://localhost:%d/%s", port, api.HealthCheckRoutePath)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
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
	var healthResponse api.HealthCheckResponse

	// Check the json can be unmarshalled
	require.NoError(t, json.Unmarshal(body, &healthResponse))

	// Check the response body
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotZero(t, healthResponse.Timestamp)
}
