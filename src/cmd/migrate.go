package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mviner000/eyygo/src/config"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func migrate() {
	dbURL := config.GetDatabaseURL()

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	// Ensure migrations table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY, filename TEXT UNIQUE)`)
	if err != nil {
		log.Fatalf("Failed to ensure migrations table exists: %v\n", err)
	}

	// List of directories to search for migrations
	migrationDirs := []string{
		"eyygo/admin/migrations",
		"eyygo/contenttypes/migrations",
		"eyygo/auth/migrations",
		"eyygo/sessions/migrations",
		// Add more directories as needed
	}

	failedDirs := []string{}
	for _, dir := range migrationDirs {
		if applyMigrationsFromDir(db, dir) {
			failedDirs = append(failedDirs, dir)
		}
	}

	if len(failedDirs) > 0 {
		log.Printf("Migrations failed in the following directories:\n")
		for _, dir := range failedDirs {
			log.Printf("- %s\n", dir)
		}
	} else {
		log.Println("All migrations applied successfully.")
	}
}

func applyMigrationsFromDir(db *sql.DB, migrationsDir string) (failed bool) {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Printf("Failed to read migrations directory %s: %v\n", migrationsDir, err)
		return true
	}

	// Sort files by name to ensure they're applied in order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			// Construct the full path of the migration file
			filePath := filepath.Join(migrationsDir, file.Name())

			// Check if this migration has already been applied
			var exists int
			err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = ?", filePath).Scan(&exists)
			if err != nil {
				log.Fatalf("Failed to check if migration %s has been applied: %v\n", filePath, err)
			}

			if exists > 0 {
				fmt.Printf("\033[34mSkipping already applied migration: %s from directory %s\033[0m\n", file.Name(), migrationsDir)
				continue
			}

			// Read and execute the SQL file
			sqlContent, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Failed to read migration file %s: %v\n", filePath, err)
			}

			_, err = db.Exec(string(sqlContent))
			if err != nil {
				fmt.Printf("\033[31mFailed to execute migration %s in directory %s: %v\033[0m\n", file.Name(), migrationsDir, err)
				return true
			}

			// Record the migration as applied
			_, err = db.Exec("INSERT INTO migrations (filename) VALUES (?)", filePath)
			if err != nil {
				log.Fatalf("Failed to record migration %s: %v\n", filePath, err)
			}

			fmt.Printf("\033[32mApplied migration: %s from directory %s\033[0m\n", file.Name(), migrationsDir)
		}
	}
	return false
}
