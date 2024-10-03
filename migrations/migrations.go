package migrations

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/mviner000/eyymi/config"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type FieldChange struct {
	Name     string
	OldType  string
	NewType  string
	Nullable bool
}

type Operation interface {
	Apply(*gorm.DB) error
	Reverse(*gorm.DB) error
	Describe() string
}

type AddColumnOperation struct {
	Table  string
	Column FieldChange
}

func (o AddColumnOperation) Apply(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", o.Table, o.Column.Name, o.Column.NewType)).Error
}

func (o AddColumnOperation) Reverse(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", o.Table, o.Column.Name)).Error
}

func (o AddColumnOperation) Describe() string {
	return fmt.Sprintf("Add column %s to %s", o.Column.Name, o.Table)
}

type DropColumnOperation struct {
	Table  string
	Column FieldChange
}

func (o DropColumnOperation) Apply(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", o.Table, o.Column.Name)).Error
}

func (o DropColumnOperation) Reverse(db *gorm.DB) error {
	return fmt.Errorf("cannot automatically reverse dropped column %s in table %s", o.Column.Name, o.Table)
}

func (o DropColumnOperation) Describe() string {
	return fmt.Sprintf("Drop column %s from %s", o.Column.Name, o.Table)
}

type AlterColumnOperation struct {
	Table  string
	Column FieldChange
}

func (o AlterColumnOperation) Apply(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s", o.Table, o.Column.Name, o.Column.NewType)).Error
}

func (o AlterColumnOperation) Reverse(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s", o.Table, o.Column.Name, o.Column.OldType)).Error
}

func (o AlterColumnOperation) Describe() string {
	return fmt.Sprintf("Alter column %s in %s", o.Column.Name, o.Table)
}

type Migration struct {
	ID           string
	Name         string
	Operations   []Operation
	Dependencies []string
}

func DetectChanges(db *gorm.DB, models ...interface{}) ([]Migration, error) {
	var migrations []Migration

	currentSchema, err := getCurrentSchemaFromDatabase(db)
	if err != nil {
		return nil, fmt.Errorf("error getting current schema from database: %w", err)
	}
	fmt.Printf("Current schema: %+v\n", currentSchema)

	for _, model := range models {
		modelType := reflect.TypeOf(model).Elem()
		modelName := modelType.Name()
		tableName := db.NamingStrategy.TableName(modelName) // Use GORM's naming strategy to get the table name

		desiredSchema, err := getDesiredSchema(db, model)
		if err != nil {
			return nil, fmt.Errorf("error getting desired schema: %w", err)
		}
		fmt.Printf("Desired schema for %s: %+v\n", tableName, desiredSchema)

		diff := compareSchemas(currentSchema[tableName], desiredSchema)
		if len(diff) > 0 {
			migration := Migration{
				ID:         generateMigrationID(),
				Name:       generateMigrationName(modelName, diff),
				Operations: diff,
			}
			migrations = append(migrations, migration)
		}
	}

	return migrations, nil
}

func getCurrentSchemaFromDatabase(db *gorm.DB) (map[string]map[string]*schema.Field, error) {
	currentSchema := make(map[string]map[string]*schema.Field)

	// Get all table names
	var tables []string
	if err := db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables).Error; err != nil {
		return nil, fmt.Errorf("error fetching table names: %w", err)
	}

	for _, tableName := range tables {
		currentSchema[tableName] = make(map[string]*schema.Field)

		// Get column information for each table
		var columns []struct {
			Name     string
			Type     string
			NotNull  int
			PK       int
			Autoincr int
		}
		if err := db.Raw(fmt.Sprintf("PRAGMA table_info('%s')", tableName)).Scan(&columns).Error; err != nil {
			return nil, fmt.Errorf("error fetching column info for table %s: %w", tableName, err)
		}

		for _, col := range columns {
			field := &schema.Field{
				Name:          col.Name,
				DBName:        col.Name,
				FieldType:     reflect.TypeOf(""), // Use string as a placeholder
				NotNull:       col.NotNull == 1,
				PrimaryKey:    col.PK == 1,
				AutoIncrement: col.Autoincr == 1,
				DataType:      getDataTypeFromSQLite(col.Type),
			}
			currentSchema[tableName][col.Name] = field
		}
	}

	return currentSchema, nil
}

func getDataTypeFromSQLite(sqliteType string) schema.DataType {
	sqliteType = strings.ToUpper(sqliteType)
	switch {
	case strings.Contains(sqliteType, "INT"):
		return schema.Int
	case strings.Contains(sqliteType, "CHAR"), strings.Contains(sqliteType, "CLOB"), strings.Contains(sqliteType, "TEXT"):
		return schema.String
	case strings.Contains(sqliteType, "REAL"), strings.Contains(sqliteType, "FLOA"), strings.Contains(sqliteType, "DOUB"):
		return schema.Float
	case strings.Contains(sqliteType, "BOOL"):
		return schema.Bool
	case strings.Contains(sqliteType, "BLOB"):
		return schema.Bytes
	case strings.Contains(sqliteType, "DATE"), strings.Contains(sqliteType, "TIME"):
		return schema.Time
	default:
		return schema.String
	}
}

func getDesiredSchema(db *gorm.DB, model interface{}) (map[string]*schema.Field, error) {
	schema := make(map[string]*schema.Field)

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, fmt.Errorf("error parsing model: %w", err)
	}

	for _, field := range stmt.Schema.Fields {
		schema[field.DBName] = field
	}

	return schema, nil
}

func compareSchemas(current, desired map[string]*schema.Field) []Operation {
	var operations []Operation

	// If either schema is nil, handle gracefully
	if current == nil && desired == nil {
		return operations
	}

	tableName := ""
	if len(desired) > 0 {
		for _, field := range desired {
			if field != nil && field.Schema != nil {
				tableName = field.Schema.Table
				break
			}
		}
	}

	if tableName == "" && len(current) > 0 {
		for _, field := range current {
			if field != nil && field.Schema != nil {
				tableName = field.Schema.Table
				break
			}
		}
	}

	if tableName == "" {
		fmt.Println("Warning: Could not determine table name")
		return operations
	}

	for field, desiredField := range desired {
		if currentField, exists := current[field]; !exists {
			operations = append(operations, AddColumnOperation{
				Table: tableName,
				Column: FieldChange{
					Name:    field,
					NewType: dataTypeToString(desiredField.DataType),
				},
			})
		} else if !fieldsEqual(currentField, desiredField) {
			operations = append(operations, AlterColumnOperation{
				Table: tableName,
				Column: FieldChange{
					Name:    field,
					OldType: dataTypeToString(currentField.DataType),
					NewType: dataTypeToString(desiredField.DataType),
				},
			})
		}
	}

	for field := range current {
		if _, exists := desired[field]; !exists {
			operations = append(operations, DropColumnOperation{
				Table: tableName,
				Column: FieldChange{
					Name: field,
				},
			})
		}
	}

	return operations
}

func fieldsEqual(a, b *schema.Field) bool {
	return a.DataType == b.DataType &&
		a.NotNull == b.NotNull &&
		a.PrimaryKey == b.PrimaryKey &&
		a.AutoIncrement == b.AutoIncrement &&
		a.Unique == b.Unique // Compare unique constraints as well
}

func dataTypeToString(dataType schema.DataType) string {
	switch dataType {
	case schema.Bool:
		return "BOOLEAN"
	case schema.Int:
		return "INTEGER"
	case schema.Uint:
		return "INTEGER UNSIGNED"
	case schema.Float:
		return "FLOAT"
	case schema.String:
		return "VARCHAR(255)"
	case schema.Time:
		return "DATETIME"
	case schema.Bytes:
		return "BLOB"
	default:
		return "TEXT"
	}
}

func generateMigrationID() string {
	return time.Now().Format("20060102150405")
}

func generateMigrationName(modelName string, operations []Operation) string {
	var parts []string
	opCounts := make(map[string]int)

	for _, op := range operations {
		opType := reflect.TypeOf(op).Name()
		opCounts[opType]++
	}

	for opType, count := range opCounts {
		switch opType {
		case "AddColumnOperation":
			parts = append(parts, fmt.Sprintf("add_%d_fields_to_%s", count, strings.ToLower(modelName)))
		case "DropColumnOperation":
			parts = append(parts, fmt.Sprintf("remove_%d_fields_from_%s", count, strings.ToLower(modelName)))
		case "AlterColumnOperation":
			parts = append(parts, fmt.Sprintf("alter_%d_fields_in_%s", count, strings.ToLower(modelName)))
		}
	}

	return strings.Join(parts, "_and_")
}

func GenerateMigration(migrations []Migration) error {
	if len(migrations) == 0 {
		fmt.Println("No changes detected.")
		return nil
	}

	for _, migration := range migrations {
		filename := fmt.Sprintf("%s_%s.go", migration.ID, migration.Name)
		migrationDir := filepath.Join(config.GetProjectRoot(), "db", "migrations")
		migrationPath := filepath.Join(migrationDir, filename)

		content := generateMigrationFile(migration)

		if err := os.MkdirAll(migrationDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating migrations directory: %w", err)
		}

		if err := ioutil.WriteFile(migrationPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("error writing migration file: %w", err)
		}

		fmt.Printf("Migration file created: %s\n", migrationPath)
	}

	return nil
}

func generateMigrationFile(migration Migration) string {
	var b strings.Builder

	writeString := func(s string) {
		_, err := b.WriteString(s)
		if err != nil {
			panic(fmt.Sprintf("Error writing string: %v", err))
		}
	}

	writeString(fmt.Sprintf(`package migrations

import (
	"gorm.io/gorm"
)

type %s struct{}

func (m *%s) Up(db *gorm.DB) error {
`, migration.Name, migration.Name))

	for _, op := range migration.Operations {
		writeString(fmt.Sprintf("\t// %s\n", op.Describe()))
		writeString("\tif err := db.Exec(\"your SQL here\").Error; err != nil {\n")
		writeString("\t\treturn err\n")
		writeString("\t}\n\n")
	}

	writeString("\treturn nil\n}\n\n")

	writeString(fmt.Sprintf(`func (m *%s) Down(db *gorm.DB) error {
`, migration.Name))

	for i := len(migration.Operations) - 1; i >= 0; i-- {
		op := migration.Operations[i]
		writeString(fmt.Sprintf("\t// Reverse: %s\n", op.Describe()))
		writeString("\tif err := db.Exec(\"your reverse SQL here\").Error; err != nil {\n")
		writeString("\t\treturn err\n")
		writeString("\t}\n\n")
	}

	writeString("\treturn nil\n}")

	return b.String()
}
