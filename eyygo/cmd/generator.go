package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mviner000/eyymi/eyygo/germ"
	"github.com/mviner000/eyymi/eyygo/germ/schema"
)

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

func generateCreateTableSQL(db *germ.DB, model interface{}) (string, error) {
	stmt := &germ.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}

	var fields []string
	for _, field := range stmt.Schema.Fields {
		fieldStr := fmt.Sprintf("%s %s", field.DBName, getDataType(field))
		if field.PrimaryKey {
			fieldStr += getPrimaryKeyString(db)
		}
		if field.NotNull {
			fieldStr += " NOT NULL"
		}
		fields = append(fields, fieldStr)
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n);", stmt.Table, strings.Join(fields, ",\n  ")), nil
}

func generateAlterTableSQL(db *germ.DB, model interface{}) ([]string, error) {
	stmt := &germ.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}

	var statements []string
	migrator := db.Migrator()

	for _, field := range stmt.Schema.Fields {
		if !migrator.HasColumn(stmt.Table, field.DBName) {
			s := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", stmt.Table, field.DBName, string(field.DataType))
			if field.NotNull {
				s += " NOT NULL"
			}
			statements = append(statements, s)
		}
	}

	return statements, nil
}

func generateUpdatedAtTriggerSQL(tableName string) []string {
	return []string{
		"-- +migrate StatementBegin",
		`CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;`,
		"-- +migrate StatementEnd",
		"",
		fmt.Sprintf(`CREATE TRIGGER update_%s_updated_at
BEFORE UPDATE ON %s
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();`, tableName, tableName),
	}
}

func getDataType(field *schema.Field) string {
	switch field.DataType {
	case "uint", "int", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
		return "INTEGER"
	case "float32", "float64":
		return "REAL"
	case "bool":
		return "BOOLEAN"
	case "string":
		return "TEXT"
	case "time.Time":
		return "DATETIME"
	default:
		// For custom types, try to determine the underlying type
		if field.FieldType != nil {
			switch field.FieldType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return "INTEGER"
			case reflect.Float32, reflect.Float64:
				return "REAL"
			case reflect.Bool:
				return "BOOLEAN"
			case reflect.String:
				return "TEXT"
			}
		}
		return "TEXT" // Default to TEXT for unknown types
	}
}

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
