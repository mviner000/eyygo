package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	conf "github.com/mviner000/eyymi/eyygo"
	"github.com/mviner000/eyymi/eyygo/config"
	"github.com/mviner000/eyymi/eyygo/germ"
	"github.com/mviner000/eyymi/eyygo/germ/driver/sqlite"
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
			log.Fatalf("Unsupported database engine: %s", conf.GetSettings().GetDatabaseConfig().Engine)
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

// Helper function to get migration files
func getMigrationFiles() ([]string, error) {
	migrationsDir := filepath.Join(conf.GetFullProjectName(), "models", "migrations")
	return filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
}

func applyMigrations(db *germ.DB) error {
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return err
	}

	for _, file := range migrationFiles {
		log.Printf("Applying migration: %s", file)

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		// Debug: Print entire content of migration file
		log.Printf("Migration file content:\n%s", string(content))

		upStatements, _, err := parseMigrationFile(string(content))
		if err != nil {
			return err
		}

		for _, stmt := range upStatements {
			log.Printf("Executing statement:\n%s", stmt)
			if err := db.Exec(stmt).Error; err != nil {
				return fmt.Errorf("error executing statement: %v\nStatement: %s", err, stmt)
			}
		}
	}

	return nil
}

func rollbackMigrations(db *germ.DB, steps int) error {
	migrationFiles, err := getMigrationFiles()
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

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		// Debug: Print entire content of migration file
		log.Printf("Migration file content:\n%s", string(content))

		_, downStatements, err := parseMigrationFile(string(content))
		if err != nil {
			return err
		}

		for _, stmt := range downStatements {
			log.Printf("Executing statement:\n%s", stmt)
			if err := db.Exec(stmt).Error; err != nil {
				return fmt.Errorf("error executing statement: %v\nStatement: %s", err, stmt)
			}
		}
	}

	return nil
}

func parseMigrationFile(content string) ([]string, []string, error) {
	var upStatements, downStatements []string
	var currentSection string
	var currentStatement strings.Builder

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "-- +migrate Up" {
			currentSection = "up"
		} else if trimmedLine == "-- +migrate Down" {
			currentSection = "down"
		} else if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "--") {
			currentStatement.WriteString(line + "\n")
			if strings.HasSuffix(trimmedLine, ";") {
				statement := strings.TrimSpace(currentStatement.String())
				if currentSection == "up" {
					upStatements = append(upStatements, statement)
				} else if currentSection == "down" {
					downStatements = append(downStatements, statement)
				}
				currentStatement.Reset()
			}
		}
	}

	// Handle any remaining statement without a semicolon
	if currentStatement.Len() > 0 {
		statement := strings.TrimSpace(currentStatement.String())
		if currentSection == "up" {
			upStatements = append(upStatements, statement)
		} else if currentSection == "down" {
			downStatements = append(downStatements, statement)
		}
	}

	return upStatements, downStatements, nil
}
