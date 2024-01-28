package cmd_test

import (
	"bytes"
	"testing"

	"github.com/jpcercal/counter/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRoot tests the root command.
func TestRoot(t *testing.T) {
	actual := new(bytes.Buffer)

	// Set the output and error streams
	cmd.RootCmd.SetOut(actual)
	cmd.RootCmd.SetErr(actual)
	err := cmd.RootCmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, actual.String(), cmd.RootCommandName)
	assert.Contains(t, actual.String(), cmd.RootCommandShort)
	assert.Contains(t, actual.String(), cmd.RootCommandLong)
}
