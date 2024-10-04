// cmd/makemigrations.go
package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/admin"
	"github.com/mviner000/eyymi/eyygo/migrations"
	"github.com/mviner000/eyymi/eyygo/operations"
	"github.com/spf13/cobra"
)

var MakeMigrationsCmd = &cobra.Command{
	Use:   "makemigrations",
	Short: "Detect model changes and create new migrations",
	Run: func(cmd *cobra.Command, args []string) {
		makeMigrations()
	},
}

func makeMigrations() {
	db, err := sql.Open("sqlite3", config.GetDatabaseURL())
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	models := admin.GetModels()

	var allOps []operations.Operation
	for _, model := range models {
		ops, err := migrations.DetectChanges(model, db)
		if err != nil {
			fmt.Printf("Error detecting changes for model %s: %v\n", model.TableName, err)
			return
		}
		allOps = append(allOps, ops...)
	}

	if len(allOps) == 0 {
		fmt.Println("No changes detected.")
		return
	}

	err = migrations.GenerateMigration(allOps)
	if err != nil {
		fmt.Printf("Error generating migration: %v\n", err)
		return
	}
}
