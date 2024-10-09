package cmd

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mviner000/eyymi/eyygo/germ"
	"github.com/mviner000/eyymi/eyygo/germ/schema"
)

// GenerateMigration generates migration SQL for provided models
func GenerateMigration(db *germ.DB, dst ...interface{}) (string, error) {
	// Log the database type
	log.Printf("Using database: %s", db.Name())

	var upStatements []string
	tableNames := make(map[string]bool)

	validate := validator.New() // Ensure the validator is initialized

	for _, model := range dst {
		if model == nil {
			return "", fmt.Errorf("model cannot be nil")
		}

		stmt := &germ.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			fmt.Println("Failed to parse model:", err, "model:", model)
			return "", err
		}

		// Validate model
		if err := validate.Struct(model); err != nil {
			return "", fmt.Errorf("validation error: %v", err)
		}

		// Generate "Up" migration
		createTableSQL, err := generateCreateTableSQL(db, stmt)
		if err != nil {
			fmt.Println("Failed to generate create table SQL:", err, "table:", stmt.Table)
			return "", err
		}

		upStatements = append(upStatements, createTableSQL)
		tableNames[stmt.Table] = true

		// Handle join tables
		joinTableSQL, err := generateJoinTableSQL(stmt)
		if err != nil {
			fmt.Println("Failed to generate join table SQL:", err, "table:", stmt.Table)
			return "", err
		}
		if joinTableSQL != "" {
			upStatements = append(upStatements, joinTableSQL)
			// Extract join table names and add them to tableNames
			joinTableNames := extractJoinTableNames(joinTableSQL)
			for _, name := range joinTableNames {
				tableNames[name] = true
			}
		}
	}

	// Remove duplicate table definitions
	uniqueUpStatements := removeDuplicateTableDefinitions(strings.Join(upStatements, "\n\n"))

	// Generate "Down" migration
	downStatements := generateDownMigration(tableNames)

	// Combine "Up" and "Down" migrations
	migration := fmt.Sprintf("-- +migrate Up\n%s\n\n-- +migrate Down\n%s",
		uniqueUpStatements,
		strings.Join(downStatements, "\n"))

	return migration, nil
}

// extractJoinTableNames extracts join table names from SQL
func extractJoinTableNames(sql string) []string {
	re := regexp.MustCompile(`CREATE TABLE IF NOT EXISTS (\w+)`)
	matches := re.FindAllStringSubmatch(sql, -1)
	var names []string
	for _, match := range matches {
		if len(match) > 1 {
			names = append(names, match[1])
		}
	}
	return names
}

// generateDownMigration generates the "Down" migration
func generateDownMigration(tableNames map[string]bool) []string {
	var downStatements []string
	for tableName := range tableNames {
		downStatements = append(downStatements, fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName))
	}
	// Sort the statements to ensure consistent output
	sort.Strings(downStatements)
	return downStatements
}

// removeDuplicateTableDefinitions removes duplicate table definitions from the SQL string
func removeDuplicateTableDefinitions(sqlString string) string {
	re := regexp.MustCompile(`(?i)CREATE TABLE IF NOT EXISTS (\w+) \((?s).*?\);`)
	matches := re.FindAllStringSubmatch(sqlString, -1)

	tables := make(map[string]string)
	var result []string

	for _, match := range matches {
		tableName := match[1]
		tableDefinition := match[0]

		if _, exists := tables[tableName]; !exists {
			tables[tableName] = tableDefinition
			result = append(result, tableDefinition)
		}
	}

	return strings.Join(result, "\n\n")
}

// generateCreateTableSQL generates the SQL for creating a table, including foreign key constraints
func generateCreateTableSQL(db *germ.DB, stmt *germ.Statement) (string, error) {
	fields := extractFields(db, stmt)

	var fieldDefs []string
	var constraints []string

	for _, field := range fields {
		if strings.Contains(field, "CONSTRAINT") {
			constraints = append(constraints, field)
		} else {
			fieldDefs = append(fieldDefs, field)
		}
	}

	// Join field definitions and constraints separately
	fieldDefsSQL := strings.Join(fieldDefs, ",\n ")
	constraintsSQL := strings.Join(constraints, ",\n ")

	// Construct the full CREATE TABLE SQL statement
	var createTableSQL string
	if constraintsSQL != "" {
		createTableSQL = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s,\n %s\n);",
			stmt.Table, fieldDefsSQL, constraintsSQL)
	} else {
		createTableSQL = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s\n);",
			stmt.Table, fieldDefsSQL)
	}

	return createTableSQL, nil
}

// extractFields generates column definitions, including foreign key, unique, and check constraints
func extractFields(db *germ.DB, stmt *germ.Statement) []string {
	var fields []string

	fmt.Printf("Extracting fields for table: %s\n", stmt.Table)

	for _, field := range stmt.Schema.Fields {
		fieldSQL := db.Migrator().FullDataTypeOf(field).SQL
		if fieldSQL != "" {
			fieldStr := fmt.Sprintf("%s %s", field.DBName, fieldSQL)

			// Handle AUTOINCREMENT for SQLite and AUTO_INCREMENT for MySQL
			if field.AutoIncrement {
				if db.Name() == "sqlite" || db.Name() == "sqlite3" {
					fieldStr = fmt.Sprintf("%s INTEGER PRIMARY KEY AUTOINCREMENT", field.DBName)
				} else if db.Name() == "mysql" {
					fieldStr = fmt.Sprintf("%s INT PRIMARY KEY AUTO_INCREMENT", field.DBName)
				}
			}

			fmt.Printf("Processing field: %s\n", field.DBName)

			// Ensure NOT NULL is only added once
			if field.NotNull {
				fieldStr = removeNotNull(fieldStr)
				fieldStr += " NOT NULL"
			}

			// Handle datetime fields
			if strings.Contains(fieldSQL, "datetime") {
				if db.Name() == "sqlite" || db.Name() == "sqlite3" {
					fieldStr = fmt.Sprintf("%s TEXT", field.DBName)
				}
			} else if db.Name() == "mysql" {
				// Convert TEXT to VARCHAR for MySQL if needed
				if strings.Contains(fieldSQL, "TEXT") {
					fieldStr = fmt.Sprintf("%s VARCHAR(255)", field.DBName)
				}
			}

			fields = append(fields, fieldStr)

			// Handle foreign key constraints
			if field.TagSettings["FOREIGNKEY"] != "" {
				refTable, refField := parseTagSetting(field.TagSettings["FOREIGNKEY"])
				if refTable == "" {
					refTable = strings.TrimSuffix(field.DBName, "ID")
					if refTable == "" {
						refTable = field.DBName
					}
				}
				if refField == "" {
					refField = "ID"
				}

				fmt.Printf("Foreign key detected for %s.%s referencing table: %s\n", stmt.Table, field.DBName, refTable)

				fkName := fmt.Sprintf("fk_%s_%s", stmt.Table, field.DBName)
				fk := fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE CASCADE",
					fkName, field.DBName, refTable, refField)
				fields = append(fields, fk)
			}

			// Handle unique constraints
			if field.Unique {
				fmt.Printf("Unique constraint detected for %s.%s\n", stmt.Table, field.DBName)
				uniName := fmt.Sprintf("uk_%s_%s", stmt.Table, field.DBName)
				uni := fmt.Sprintf("CONSTRAINT %s UNIQUE (%s)", uniName, field.DBName)
				fields = append(fields, uni)
			}
		}
	}

	return fields
}

// parseTagSetting parses the FOREIGNKEY tag setting
func parseTagSetting(tagSetting string) (string, string) {
	parts := strings.Split(tagSetting, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

// hasForeignKeyConstraint checks if a foreign key constraint is defined for the field
func hasForeignKeyConstraint(field *schema.Field) (bool, string) {
	for key, value := range field.TagSettings {
		if strings.EqualFold(key, "foreignKey") {
			fmt.Printf("Foreign key tag detected for field %s: %s\n", field.Name, value)
			return true, value
		}
	}
	return false, "id" // Default to 'id' if no specific field is specified
}

// removeNotNull removes existing NOT NULL constraint
func removeNotNull(fieldStr string) string {
	return strings.ReplaceAll(fieldStr, " NOT NULL", "")
}

// generateJoinTableSQL creates SQL for join tables, including foreign key constraints
func generateJoinTableSQL(stmt *germ.Statement) (string, error) {
	var joinTableSQL []string

	for _, rel := range stmt.Schema.Relationships.Relations {
		if rel.Type == schema.Many2Many && rel.JoinTable != nil {
			fmt.Printf("Generating join table for %s and %s\n", stmt.Table, rel.FieldSchema.Table)
			joinTable := rel.JoinTable

			// Use the correct field to get the reference table
			refTable := rel.FieldSchema.Table // Example, adjust according to your actual struct

			joinFields := []string{
				fmt.Sprintf("%s_id INTEGER NOT NULL", stmt.Table),
				fmt.Sprintf("%s_id INTEGER NOT NULL", rel.FieldSchema.Table),
				fmt.Sprintf("PRIMARY KEY (%s_id, %s_id)", stmt.Table, rel.FieldSchema.Table),
				fmt.Sprintf("FOREIGN KEY (%s_id) REFERENCES %s(id) ON DELETE CASCADE", stmt.Table, refTable),
				fmt.Sprintf("FOREIGN KEY (%s_id) REFERENCES %s(id) ON DELETE CASCADE", rel.FieldSchema.Table, refTable),
			}

			joinTableSQL = append(joinTableSQL, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s\n);",
				joinTable.Name, strings.Join(joinFields, ",\n ")))
		}
	}

	return strings.Join(joinTableSQL, "\n\n"), nil
}
