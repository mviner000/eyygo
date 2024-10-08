package cmd

import (
	"log"
	"os"

	"github.com/mviner000/eyymi/eyygo/core"
	models "github.com/mviner000/eyymi/eyygo/post"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ServerCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Start the development server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Initializing database connection...")
		// Initialize database
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("GERM TEST DB Failed: Unable to connect to database: %v", err)
		}
		log.Println("GERM Database connection established successfully.")

		log.Println("Starting auto-migration for Post model...")
		// Perform auto migration
		if err := db.AutoMigrate(&models.Post{}); err != nil {
			log.Fatalf("GERM Auto-Migration Failed: Unable to migrate Post model: %v", err)
		}
		log.Println("Auto-migration completed successfully.")

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
