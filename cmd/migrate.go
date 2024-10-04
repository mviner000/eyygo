// cmd/migrate.go
package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

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

	migrationsDir := "admin/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v\n", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(migrationsDir, file.Name())
			sqlContent, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Failed to read migration file %s: %v\n", file.Name(), err)
			}

			_, err = db.Exec(string(sqlContent))
			if err != nil {
				log.Fatalf("Failed to execute migration %s: %v\n", file.Name(), err)
			}

			fmt.Printf("Applied migration: %s\n", file.Name())
		}
	}
}
