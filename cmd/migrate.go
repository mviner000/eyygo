package cmd

import (
	"log"
	"time"

	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/types"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

func runMigrations() {
	dbURL := config.GetDatabaseURL()
	config.DebugLog("Using database URL: %s", dbURL)

	// Open a connection to the database
	db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&types.User{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully.")

	// Create a superuser if it doesn't exist
	createSuperUserIfNotExists(db)
}

func createSuperUserIfNotExists(db *gorm.DB) {
	var count int64
	db.Model(&types.User{}).Where("is_superuser = ?", true).Count(&count)
	if count == 0 {
		superuser := types.User{
			Username:    "admin",
			Email:       "admin@example.com",
			Password:    "adminpassword", // In a real app, hash this password
			DateJoined:  time.Now(),
			IsActive:    true,
			IsStaff:     true,
			IsSuperuser: true,
		}
		result := db.Create(&superuser)
		if result.Error != nil {
			log.Fatalf("Failed to create superuser: %v", result.Error)
		}
		log.Println("Superuser created successfully.")
	} else {
		log.Println("Superuser already exists.")
	}
}
