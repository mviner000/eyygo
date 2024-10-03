package cmd

import (
	"fmt"

	"github.com/mviner000/eyymi/admin"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/migrations"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var MakeMigrationsCmd = &cobra.Command{
	Use:   "makemigrations",
	Short: "Detect model changes and create new migrations",
	Run: func(cmd *cobra.Command, args []string) {
		makeMigrations()
	},
}

func makeMigrations() {
	db, err := gorm.Open(sqlite.Open(config.GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	// Add all your models here
	models := []interface{}{
		&admin.User{},
		// Add other models here
	}

	diffs, err := migrations.DetectChanges(db, models...)
	if err != nil {
		fmt.Printf("Error detecting changes: %v\n", err)
		return
	}

	err = migrations.GenerateMigration(diffs)
	if err != nil {
		fmt.Printf("Error generating migration: %v\n", err)
		return
	}
}
