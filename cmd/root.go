package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	RootCommandName  = "counter"
	RootCommandShort = "A simple counter api"
	RootCommandLong  = `A simple counter api with support to concurrency and graceful shutdown`
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   RootCommandName,
	Short: RootCommandShort,
	Long:  RootCommandLong,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Printf("error executing root command: %v\n", err)
		os.Exit(1)
	}
}

// init initializes cobra.
func init() {
	cobra.OnInitialize()
}
