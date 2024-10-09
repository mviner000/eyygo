package cmd

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/germ"
	"github.com/mviner000/eyymi/eyygo/germ/driver/sqlite"
	"github.com/mviner000/eyymi/project_name"
	"github.com/spf13/cobra"
)

var (
	rollback bool
	steps    int
)

var MigratorCmd = &cobra.Command{
	Use:   "migrator",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Initializing database connection for migration...")

		// Get database URL
		dbURL := config.GetDatabaseURL()
		if dbURL == "" {
			log.Fatalf("Unsupported database engine: %s", project_name.AppSettings.GetDatabaseConfig().Engine)
		}

		log.Printf("Using database: %s", dbURL)

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

		log.Println("GERM Database connection established successfully.")

		if rollback {
			err = rollbackMigrations(db, steps)
		} else {
			err = applyMigrations(db)
		}

		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		log.Println("Migration completed successfully.")
	},
}

func init() {
	MigratorCmd.Flags().BoolVarP(&rollback, "rollback", "r", false, "Rollback migrations")
	MigratorCmd.Flags().IntVarP(&steps, "steps", "s", 1, "Number of migrations to roll back")
}

func applyMigrations(db *germ.DB) error {
	migrationsDir := "migrations"
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}

	for _, file := range migrationFiles {
		log.Printf("Applying migration: %s", file)

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		upStatements, _, err := parseMigrationFile(string(content))
		if err != nil {
			return err
		}

		for _, stmt := range upStatements {
			if err := db.Exec(stmt).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func rollbackMigrations(db *germ.DB, steps int) error {
	migrationsDir := "migrations"
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}

	// Sort migration files in reverse order
	for i := len(migrationFiles)/2 - 1; i >= 0; i-- {
		opp := len(migrationFiles) - 1 - i
		migrationFiles[i], migrationFiles[opp] = migrationFiles[opp], migrationFiles[i]
	}

	for i, file := range migrationFiles {
		if i >= steps {
			break
		}

		log.Printf("Rolling back migration: %s", file)

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, downStatements, err := parseMigrationFile(string(content))
		if err != nil {
			return err
		}

		for _, stmt := range downStatements {
			if err := db.Exec(stmt).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func parseMigrationFile(content string) ([]string, []string, error) {
	var upStatements, downStatements []string
	var currentSection string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "-- +migrate Up" {
			currentSection = "up"
		} else if line == "-- +migrate Down" {
			currentSection = "down"
		} else if line != "" && !strings.HasPrefix(line, "--") {
			if currentSection == "up" {
				upStatements = append(upStatements, line)
			} else if currentSection == "down" {
				downStatements = append(downStatements, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return upStatements, downStatements, nil
}
