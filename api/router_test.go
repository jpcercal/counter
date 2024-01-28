package api_test

import (
	"testing"

	"github.com/jpcercal/counter/api"
	"github.com/stretchr/testify/require"
)

// TestRouter tests the router.
func TestRouter(t *testing.T) {
	_, err := api.NewHTTPHandler()
	require.NoError(t, err)
}
