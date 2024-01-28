package cmd

import (
	"log"

	"github.com/jpcercal/counter/api"
	"github.com/jpcercal/counter/config"
	"github.com/spf13/cobra"
)

// httpServerPort is the port used by the http server.
const httpServerPort = config.ServerHTTPPort

const (
	ServeCommandName  = "serve"
	ServeCommandShort = "Starts a http server and serves the configured api"
	ServeCommandLong  = `Starts a http server and serves the configured api`
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   ServeCommandName,
	Short: ServeCommandShort,
	Long:  ServeCommandLong,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := api.NewServer(httpServerPort)
		if err != nil {
			log.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
