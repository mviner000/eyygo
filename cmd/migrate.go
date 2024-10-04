// cmd/migrate.go
package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mviner000/eyymi/config"
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
	db, err := sql.Open("sqlite3", config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	// Ensure migrations table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY, filename TEXT UNIQUE)`)
	if err != nil {
		log.Fatalf("Failed to ensure migrations table exists: %v\n", err)
	}

	migrationsDir := "admin/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v\n", err)
	}

	// Sort files by name to ensure they're applied in order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			// Check if this migration has already been applied
			var exists int
			err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = ?", file.Name()).Scan(&exists)
			if err != nil {
				log.Fatalf("Failed to check if migration %s has been applied: %v\n", file.Name(), err)
			}

			if exists > 0 {
				fmt.Printf("Skipping already applied migration: %s\n", file.Name())
				continue
			}

			// Read and execute the SQL file
			filePath := filepath.Join(migrationsDir, file.Name())
			sqlContent, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Failed to read migration file %s: %v\n", file.Name(), err)
			}

			_, err = db.Exec(string(sqlContent))
			if err != nil {
				log.Fatalf("Failed to execute migration %s: %v\n", file.Name(), err)
			}

			// Record the migration as applied
			_, err = db.Exec("INSERT INTO migrations (filename) VALUES (?)", file.Name())
			if err != nil {
				log.Fatalf("Failed to record migration %s: %v\n", file.Name(), err)
			}

			fmt.Printf("Applied migration: %s\n", file.Name())
		}
	}
}
