package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tanq16/walkthrough/internal/server"
	u "github.com/tanq16/walkthrough/internal/utils"
)

var serveFlags struct {
	port int
	host string
	data string
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the canvas viewer web server",
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.New(serveFlags.port, serveFlags.host, serveFlags.data)
		if err := srv.Setup(); err != nil {
			u.PrintFatal("Failed to setup server", err)
		}
		u.PrintInfo(fmt.Sprintf("Starting server on %s:%d", serveFlags.host, serveFlags.port))
		u.PrintInfo(fmt.Sprintf("Serving canvas from: %s", serveFlags.data))
		if err := srv.Run(); err != nil {
			u.PrintFatal("Server error", err)
		}
	},
}

func init() {
	serveCmd.Flags().IntVarP(&serveFlags.port, "port", "p", 8080, "Port to listen on")
	serveCmd.Flags().StringVarP(&serveFlags.host, "host", "H", "0.0.0.0", "Host to bind to")
	serveCmd.Flags().StringVarP(&serveFlags.data, "data", "d", ".", "Directory containing data.json and local files")
}
