package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/mviner000/eyymi/eyygo/germ"
	"github.com/mviner000/eyymi/eyygo/germ/schema"
)

// GenerateMigration generates migration SQL for provided models
func GenerateMigration(db *germ.DB, dst ...interface{}) (string, error) {
	var upStatements []string
	tableNames := make(map[string]bool)

	for _, model := range dst {
		stmt := &germ.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return "", err
		}

		// Generate "Up" migration
		createTableSQL, err := generateCreateTableSQL(db, model)
		if err != nil {
			return "", err
		}

		upStatements = append(upStatements, createTableSQL)
		tableNames[stmt.Table] = true

		// Handle join tables
		joinTableSQL, err := generateJoinTableSQL(stmt)
		if err != nil {
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

// Helper function to extract join table names from SQL
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

// Helper function to generate "Down" migration
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

// generateCreateTableSQL generates the SQL for creating a table
func generateCreateTableSQL(db *germ.DB, model interface{}) (string, error) {
	stmt := &germ.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}

	fields := extractFields(db, stmt)

	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s\n);",
		stmt.Table, strings.Join(fields, ",\n "))

	joinTableSQL, err := generateJoinTableSQL(stmt)
	if err != nil {
		return "", err
	}

	if joinTableSQL != "" {
		createTableSQL += "\n\n" + joinTableSQL
	}

	return createTableSQL, nil
}

func extractFields(db *germ.DB, stmt *germ.Statement) []string {
	var fields []string

	for _, field := range stmt.Schema.Fields {
		fieldSQL := db.Migrator().FullDataTypeOf(field).SQL
		if fieldSQL != "" {
			fieldStr := fmt.Sprintf("%s %s", field.DBName, fieldSQL)

			// Ensure NOT NULL is only added once
			if field.NotNull {
				fieldStr = removeNotNull(fieldStr) // Remove any existing NOT NULL
				fieldStr += " NOT NULL"
			}

			fields = append(fields, fieldStr)

			if foreignKey := field.TagSettings["FOREIGNKEY"]; foreignKey != "" {
				refTable := field.TagSettings["REFERENCES"]
				if refTable == "" {
					refTable = foreignKey
				}
				fkName := fmt.Sprintf("fk_%s_%s", stmt.Table, field.DBName)
				fk := fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(id)",
					fkName, field.DBName, refTable)
				fields = append(fields, fk)
			}
		}
	}

	return fields
}

// Helper function to remove existing NOT NULL constraint
func removeNotNull(fieldStr string) string {
	return strings.ReplaceAll(fieldStr, " NOT NULL", "")
}

func generateJoinTableSQL(stmt *germ.Statement) (string, error) {
	var joinTableSQL []string

	for _, rel := range stmt.Schema.Relationships.Relations {
		if rel.Type == schema.Many2Many && rel.JoinTable != nil {
			joinTable := rel.JoinTable
			joinFields := []string{
				fmt.Sprintf("%s_id INTEGER NOT NULL", stmt.Table),
				fmt.Sprintf("%s_id INTEGER NOT NULL", rel.FieldSchema.Table),
				fmt.Sprintf("PRIMARY KEY (%s_id, %s_id)", stmt.Table, rel.FieldSchema.Table),
				fmt.Sprintf("FOREIGN KEY (%s_id) REFERENCES %s(id)", stmt.Table, stmt.Table),
				fmt.Sprintf("FOREIGN KEY (%s_id) REFERENCES %s(id)", rel.FieldSchema.Table, rel.FieldSchema.Table),
			}
			joinTableSQL = append(joinTableSQL, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s\n);",
				joinTable.Name, strings.Join(joinFields, ",\n ")))
		}
	}

	return strings.Join(joinTableSQL, "\n\n"), nil
}
