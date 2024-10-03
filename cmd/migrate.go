package cmd

import (
	"log"
	"time"

	"github.com/mviner000/eyymi/admin"
	"github.com/mviner000/eyymi/config"
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
	db, err := gorm.Open(sqlite.Open(config.GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&admin.User{},
		&admin.Group{},
		&admin.UserGroup{},
		// Add other models here
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully.")

	// Create default groups
	createDefaultGroups(db)

	// Create a superuser if it doesn't exist
	createSuperUserIfNotExists(db)
}

func createDefaultGroups(db *gorm.DB) {
	defaultGroups := []admin.Group{
		{Name: "Administrators", Description: "Users with full access"},
		{Name: "Staff", Description: "Users with limited administrative access"},
		{Name: "Users", Description: "Regular users"},
	}

	for _, group := range defaultGroups {
		var existingGroup admin.Group
		if err := db.Where("name = ?", group.Name).First(&existingGroup).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&group).Error; err != nil {
					log.Printf("Failed to create group %s: %v", group.Name, err)
				} else {
					log.Printf("Group %s created successfully", group.Name)
				}
			} else {
				log.Printf("Error checking for existing group %s: %v", group.Name, err)
			}
		} else {
			log.Printf("Group %s already exists", group.Name)
		}
	}
}

func createSuperUserIfNotExists(db *gorm.DB) {
	var count int64
	db.Model(&admin.User{}).Where("is_superuser = ?", true).Count(&count)
	if count == 0 {
		adminGroup := admin.Group{}
		err := db.Where("name = ?", "Administrators").First(&adminGroup).Error
		if err != nil {
			log.Printf("Failed to find Administrators group: %v", err)
		}

		superuser := admin.User{
			Username:    "admin",
			Email:       "admin@example.com",
			Password:    "adminpassword", // In a real app, hash this password
			DateJoined:  time.Now(),
			IsActive:    true,
			IsStaff:     true,
			IsSuperuser: true,
			Groups:      []admin.Group{adminGroup},
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
