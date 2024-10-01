package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mviner000/eyymi/config"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "manage",
	Short: "Project management tool for your Go application",
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

type User struct {
	gorm.Model
	Username    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	DateJoined  time.Time
	IsActive    bool `gorm:"default:true"`
	IsStaff     bool `gorm:"default:false"`
	IsSuperuser bool `gorm:"default:false"`
}

func runMigrations() {
	dbURL := config.GetDatabaseURL()
	log.Printf("Debug: Using database URL: %s", dbURL)

	// Open a connection to the database
	db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully.")

	// Create a superuser if it doesn't exist
	createSuperUserIfNotExists(db)
}

func createSuperUserIfNotExists(db *gorm.DB) {
	var count int64
	db.Model(&User{}).Where("is_superuser = ?", true).Count(&count)

	if count == 0 {
		superuser := User{
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
