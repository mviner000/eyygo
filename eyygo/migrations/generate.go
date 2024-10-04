package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mviner000/eyymi/eyygo/operations"
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

	// Determine the next migration number
	nextNumber, err := getNextMigrationNumber(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to determine next migration number: %v", err)
	}

	// Describe the proposed change
	changeDescription := describeChange(ops)

	// Create the migration file with a descriptive name
	filename := filepath.Join(migrationsDir, fmt.Sprintf("%d_%s.sql", nextNumber, changeDescription))
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
			_, err = file.WriteString(sql + "\n")
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

func getNextMigrationNumber(dir string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	maxNumber := 0
	re := regexp.MustCompile(`^(\d+)_`)

	for _, file := range files {
		if matches := re.FindStringSubmatch(file.Name()); matches != nil {
			number, err := strconv.Atoi(matches[1])
			if err == nil && number > maxNumber {
				maxNumber = number
			}
		}
	}

	return maxNumber + 1, nil
}

func describeChange(ops []operations.Operation) string {
	// Initialize a map to count occurrences of each operation type
	opCount := make(map[string]int)

	// Iterate over operations to generate a description
	for _, op := range ops {
		switch o := op.(type) {
		case *operations.AddTable:
			opCount[fmt.Sprintf("add_table_%s", o.Model.TableName)]++
		case *operations.AddField:
			opCount[fmt.Sprintf("add_field_%s_to_%s", o.FieldName, o.ModelName)]++
		case *operations.RemoveField:
			opCount[fmt.Sprintf("remove_field_%s_from_%s", o.FieldName, o.ModelName)]++
			// Add more cases as needed for other operations
		}
	}

	// Build a description string
	var descriptions []string
	for desc, count := range opCount {
		if count > 1 {
			descriptions = append(descriptions, fmt.Sprintf("%s_%d", desc, count))
		} else {
			descriptions = append(descriptions, desc)
		}
	}

	// Join descriptions with underscores
	return strings.Join(descriptions, "_")
}
