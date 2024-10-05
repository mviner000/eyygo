package cmd

import (
	"log"
	"os"

	"github.com/mviner000/eyymi/eyygo/core"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Start the development server",
	Run: func(cmd *cobra.Command, args []string) {
		httpPort := os.Getenv("HTTP_PORT")

		if httpPort == "" {
			httpPort = "8000"
		}

		log.Printf("Starting server on port %s", httpPort)

		// Start the server
		app := core.NewApp()
		if err := app.Listen(":" + httpPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}
