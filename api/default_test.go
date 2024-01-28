package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jpcercal/counter/api"
	testing_helper "github.com/jpcercal/counter/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultHandler tests the default handler.
func TestDefaultHandler(t *testing.T) {
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
	url := fmt.Sprintf("http://localhost:%d/%s", port, "not-found")
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	require.NoError(t, err)

	// Send the request
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	// Close the request body when test finishes
	defer response.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
