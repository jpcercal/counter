package cmd_test

import (
	"bytes"
	"testing"

	"github.com/jpcercal/counter/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServe tests the serve command.
func TestServe(t *testing.T) {
	actual := new(bytes.Buffer)

	// Set the output and error streams
	cmd.RootCmd.SetOut(actual)
	cmd.RootCmd.SetErr(actual)
	err := cmd.RootCmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, actual.String(), cmd.ServeCommandName)
	assert.Contains(t, actual.String(), cmd.ServeCommandShort)
	assert.Contains(t, actual.String(), cmd.ServeCommandLong)
}
