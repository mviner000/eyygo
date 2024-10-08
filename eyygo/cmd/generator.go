package cmd

import (
	"fmt"
	"strings"

	"github.com/mviner000/eyymi/eyygo/germ"
)

// GenerateMigration generates migration SQL for provided models
func GenerateMigration(db *germ.DB, dst ...interface{}) (string, error) {
	var upStatements, downStatements []string

	for _, model := range dst {
		stmt := &germ.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return "", err
		}

		tableName := stmt.Table

		// Generate "Up" migration
		createTableSQL, err := generateCreateTableSQL(db, model)
		if err != nil {
			return "", err
		}
		upStatements = append(upStatements, createTableSQL)

		// Generate "Down" migration
		dropTableSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)
		downStatements = append(downStatements, dropTableSQL)
	}

	// Combine "Up" and "Down" migrations
	migration := fmt.Sprintf("-- +migrate Up\n%s\n\n-- +migrate Down\n%s",
		strings.Join(upStatements, "\n\n"),
		strings.Join(downStatements, "\n"))

	return migration, nil
}

// generateCreateTableSQL generates the SQL for creating a table
func generateCreateTableSQL(db *germ.DB, model interface{}) (string, error) {
	stmt := &germ.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}

	var fields []string
	for _, field := range stmt.Schema.Fields {
		expr := db.Migrator().FullDataTypeOf(field)
		fieldType := fmt.Sprintf("%v", expr.SQL)

		// Split the fieldType into its components
		parts := strings.Fields(fieldType)
		dataType := parts[0]

		constraints := make(map[string]bool)
		for _, part := range parts[1:] {
			constraints[strings.ToUpper(part)] = true
		}

		// Add missing constraints
		if field.PrimaryKey && !constraints["PRIMARY"] {
			constraints["PRIMARY KEY AUTOINCREMENT"] = true
		}
		if field.NotNull && !constraints["NOT"] {
			constraints["NOT NULL"] = true
		}
		if field.Unique && !constraints["UNIQUE"] {
			constraints["UNIQUE"] = true
		}

		// Construct the field string
		fieldStr := fmt.Sprintf("%s %s", field.Name, dataType)
		for constraint := range constraints {
			fieldStr += " " + constraint
		}

		fields = append(fields, fieldStr)
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s\n);", stmt.Table, strings.Join(fields, ",\n ")), nil
}

// getPrimaryKeyString returns the primary key definition based on the database dialect
func getPrimaryKeyString(db *germ.DB) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return " PRIMARY KEY AUTOINCREMENT"
	case "mysql":
		return " PRIMARY KEY AUTO_INCREMENT"
	case "postgres":
		return " PRIMARY KEY"
	default:
		return " PRIMARY KEY"
	}
}
