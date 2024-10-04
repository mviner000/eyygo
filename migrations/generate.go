// migrations/generate.go
package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mviner000/eyymi/operations"
)

func GenerateMigration(ops []operations.Operation) error {
	if len(ops) == 0 {
		return fmt.Errorf("no operations to generate migration for")
	}

	// Define the directory for storing migrations
	migrationsDir := "admin/migrations"

	// Ensure the directory exists
	if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}

	// Create the migration file with a timestamp
	timestamp := time.Now().Format("20060102150405")
	filename := filepath.Join(migrationsDir, fmt.Sprintf("%s_auto_generated.sql", timestamp))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %v", err)
	}
	defer file.Close()

	for _, op := range ops {
		if sqlOp, ok := op.(interface{ SQL() (string, error) }); ok {
			sql, err := sqlOp.SQL()
			if err != nil {
				return fmt.Errorf("failed to generate SQL for operation: %v", err)
			}
			_, err = file.WriteString(sql + ";\n")
			if err != nil {
				return fmt.Errorf("failed to write SQL to migration file: %v", err)
			}
		} else {
			return fmt.Errorf("operation does not implement SQL() method: %T", op)
		}
	}

	fmt.Printf("Migration file created: %s\n", filename)
	return nil
}
