package cmd

import (
	"log"
	"os"

	"github.com/mviner000/eyygo/src/core"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Start the development server",
	Run: func(cmd *cobra.Command, args []string) {
		httpPort := os.Getenv("HTTP_PORT")
		if httpPort == "" {
			httpPort = "8000"
			log.Println("No HTTP_PORT specified in environment. Using default port 8000.")
		}
		log.Printf("Preparing to start server on port %s", httpPort)

		// Start the server
		app := core.NewApp()
		log.Printf("Server initialized. Attempting to listen on port %s...", httpPort)
		if err := app.Listen(":" + httpPort); err != nil {
			log.Fatalf("Server Startup Failed: Unable to start server: %v", err)
		}
	},
}
