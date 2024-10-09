package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/germ"
	"github.com/mviner000/eyymi/eyygo/germ/driver/sqlite"
	"github.com/mviner000/eyymi/eyygo/registry"
	"github.com/mviner000/eyymi/project_name"
	models "github.com/mviner000/eyymi/project_name/posts"

	"github.com/spf13/cobra"
)

var MakeMigrationCmd = &cobra.Command{
	Use:   "makemigrations",
	Short: "Create a new migration file",
	Run: func(cmd *cobra.Command, args []string) {

		models.RegisterModels() // Registering models

		log.Println("Creating new migration file...")

		// Get database URL
		dbURL := config.GetDatabaseURL()
		if dbURL == "" {
			log.Fatalf("Unsupported database engine: %s", project_name.AppSettings.GetDatabaseConfig().Engine)
		}

		// Initialize database
		db, err := germ.Open(sqlite.Open(dbURL), &germ.Config{})
		if err != nil {
			log.Fatalf("GERM DB Failed: Unable to connect to database: %v", err)
		}

		// Get the underlying sql.DB
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get underlying SQL DB: %v", err)
		}
		defer sqlDB.Close()

		// Generate migration content
		migrationContent, err := generateMigrationContent(db)
		if err != nil {
			log.Fatalf("Failed to generate migration content: %v", err)
		}

		// Create migration file
		filename, err := createMigrationFile(migrationContent)
		if err != nil {
			log.Fatalf("Failed to create migration file: %v", err)
		}

		log.Printf("Migrations for 'posts':\nposts/migrations/%s", filename)
		log.Println("Migration file created successfully.")
	},
}

func generateMigrationContent(db *germ.DB) (string, error) {
	generator := NewMigrationGenerator(db)

	// Retrieve all registered model names from the registry
	modelNames := registry.GetRegisteredModelNames()

	if len(modelNames) == 0 {
		log.Println("[WARN] No models registered for migration.")
		return "", fmt.Errorf("no registered models found")
	}

	// Convert []string to []interface{} for the migration generator
	var modelInterfaces []interface{}
	for _, name := range modelNames {
		// Retrieve the actual model instances using registry
		model, ok := registry.GetRegisteredModel(name) // Assuming you have this method to get the model by name
		if !ok {
			log.Printf("[WARN] Model %s not found in registry.\n", name)
			continue
		}

		modelInterfaces = append(modelInterfaces, model)
		log.Printf("[INFO] Found model for migration: %s\n", name) // Log found models
	}

	if len(modelInterfaces) == 0 {
		log.Println("[WARN] No valid models found for migration.")
		return "", fmt.Errorf("no valid models found for migration")
	}

	// Now pass modelInterfaces to the GenerateMigration method
	return generator.GenerateMigration(modelInterfaces...)
}

func createMigrationFile(content string) (string, error) {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_migration.sql", timestamp)
	migrationsDir := filepath.Join("project_name", "posts", "migrations")

	if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
		return "", err
	}

	filePath := filepath.Join(migrationsDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return filename, nil
}
