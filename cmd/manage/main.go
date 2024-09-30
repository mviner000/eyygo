package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // or whatever database you're using
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/mviner000/eyymi/config" // Update this import path to match your project structure
)

var rootCmd = &cobra.Command{
	Use:   "manage",
	Short: "Project management tool for your Go application",
}

func init() {
	rootCmd.AddCommand(startappCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(makemigrationsCmd)
}

var startappCmd = &cobra.Command{
	Use:   "startapp [name]",
	Short: "Create a new application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		createApp(appName)
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

var makemigrationsCmd = &cobra.Command{
	Use:   "makemigrations [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		createMigration(migrationName)
	},
}

func createApp(name string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}

	appDir := filepath.Join(cwd, name)
	absAppDir, err := filepath.Abs(appDir)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		return
	}

	err = os.MkdirAll(absAppDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	files := map[string]func(string) string{
		"handlers.go": createHandlersTemplate,
		"models.go":   createModelsTemplate,
		"routes.go":   createRoutesTemplate,
	}

	for file, templateFunc := range files {
		fullPath := filepath.Join(absAppDir, file)
		content := templateFunc(name)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", fullPath, err)
			return
		}
		fmt.Printf("Created %s\n", fullPath)
	}

	config.AddInstalledApp(name)

	fmt.Printf("Application '%s' created successfully and added to INSTALLED_APPS.\n", name)
	fmt.Printf("Full path of the new application: %s\n", absAppDir)
}

func createHandlersTemplate(appName string) string {
	return fmt.Sprintf(`package %s

import (
	"github.com/gofiber/fiber/v2"
)

// HelloHandler handles the hello route
func HelloHandler(c *fiber.Ctx) error {
	return c.SendString("Hello from %s!")
}
`, appName, appName)
}

func createModelsTemplate(appName string) string {
	return fmt.Sprintf(`package %s

// Example model
type Example struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`, appName)
}

func createRoutesTemplate(appName string) string {
	return fmt.Sprintf(`package %s

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures the routes for the %s app
func SetupRoutes(app *fiber.App) {
	%sGroup := app.Group("/%s")
	%sGroup.Get("/hello", HelloHandler)
}
`, appName, appName, strings.ToLower(appName), strings.ToLower(appName), strings.ToLower(appName))
}

func runMigrations() {
	// Use the database URL from your config
	dbURL := config.GetDatabaseURL()
	migrationsPath := filepath.Join("..", "..", "migrations")
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dbURL)
	if err != nil {
		fmt.Printf("Error creating migrate instance: %v\n", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Printf("Error running migrations: %v\n", err)
		return
	}

	fmt.Println("Migrations completed successfully.")
}

func createMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	upFileName := fmt.Sprintf("%s_%s.up.sql", timestamp, name)
	downFileName := fmt.Sprintf("%s_%s.down.sql", timestamp, name)

	migrationsDir := filepath.Join("..", "..", "migrations")
	os.MkdirAll(migrationsDir, os.ModePerm)

	for _, fileName := range []string{upFileName, downFileName} {
		f, err := os.Create(filepath.Join(migrationsDir, fileName))
		if err != nil {
			fmt.Printf("Error creating migration file %s: %v\n", fileName, err)
			return
		}
		defer f.Close()
		fmt.Printf("Created migration file: %s\n", fileName)
	}

	fmt.Println("Migration files created successfully. Please edit them to add your migration SQL.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	installedApps := config.GetInstalledApps()
    fmt.Println("Installed Apps:", installedApps)
}