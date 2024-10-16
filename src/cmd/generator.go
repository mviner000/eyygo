package cmd

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mviner000/eyymi/src/germ"
	"github.com/mviner000/eyymi/src/germ/schema"
)

// MigrationGenerator handles the generation of SQL migrations
type MigrationGenerator struct {
	db       *germ.DB
	validate *validator.Validate
	tables   map[string]*TableInfo
}

// TableInfo stores information about a table and its dependencies
type TableInfo struct {
	Name         string
	Dependencies []string
	SQL          string
}

// NewMigrationGenerator creates a new MigrationGenerator
func NewMigrationGenerator(db *germ.DB) *MigrationGenerator {
	return &MigrationGenerator{
		db:       db,
		validate: validator.New(),
		tables:   make(map[string]*TableInfo),
	}
}

// GenerateMigration generates migration SQL for provided models
func (mg *MigrationGenerator) GenerateMigration(dst ...interface{}) (string, error) {
	log.Printf("Using database: %s", mg.db.Name())

	for _, model := range dst {
		if model == nil {
			return "", fmt.Errorf("model cannot be nil")
		}

		if err := mg.validate.Struct(model); err != nil {
			return "", fmt.Errorf("validation error: %v", err)
		}

		stmt := &germ.Statement{DB: mg.db}
		if err := stmt.Parse(model); err != nil {
			return "", fmt.Errorf("failed to parse model: %v", err)
		}

		createTableSQL, deps, err := mg.generateCreateTableSQL(stmt)
		if err != nil {
			return "", fmt.Errorf("failed to generate create table SQL for %s: %v", stmt.Table, err)
		}

		mg.tables[stmt.Table] = &TableInfo{
			Name:         stmt.Table,
			Dependencies: deps,
			SQL:          createTableSQL,
		}

		joinTableSQL, joinDeps, err := mg.generateJoinTableSQL(stmt)
		if err != nil {
			return "", fmt.Errorf("failed to generate join table SQL for %s: %v", stmt.Table, err)
		}
		if joinTableSQL != "" {
			joinTableName := mg.extractJoinTableNames(joinTableSQL)[0]
			mg.tables[joinTableName] = &TableInfo{
				Name:         joinTableName,
				Dependencies: joinDeps,
				SQL:          joinTableSQL,
			}
		}
	}

	sortedTables, err := mg.priorityAwareSort()
	if err != nil {
		return "", fmt.Errorf("failed to sort tables: %v", err)
	}

	var upStatements []string
	for _, tableName := range sortedTables {
		upStatements = append(upStatements, mg.tables[tableName].SQL)
	}

	uniqueUpStatements := mg.removeDuplicateTableDefinitions(strings.Join(upStatements, "\n\n"))
	downStatements := mg.generateDownMigration(sortedTables)

	migration := fmt.Sprintf("-- +migrate Up\n%s\n\n-- +migrate Down\n%s",
		uniqueUpStatements,
		strings.Join(downStatements, "\n"))

	return migration, nil
}

func (mg *MigrationGenerator) priorityAwareSort() ([]string, error) {
	var noDeps, withDeps []string
	visited := make(map[string]bool)
	var visit func(string) error

	visit = func(name string) error {
		if visited[name] {
			return nil
		}
		if mg.tables[name] == nil {
			return fmt.Errorf("unknown table: %s", name)
		}
		visited[name] = true
		for _, dep := range mg.tables[name].Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}
		if len(mg.tables[name].Dependencies) == 0 {
			noDeps = append(noDeps, name)
		} else {
			withDeps = append(withDeps, name)
		}
		return nil
	}

	for name := range mg.tables {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	// Sort tables without dependencies alphabetically
	sort.Strings(noDeps)

	// Combine no-dependency tables with dependency tables
	return append(noDeps, withDeps...), nil
}

func (mg *MigrationGenerator) extractJoinTableNames(sql string) []string {
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

func (mg *MigrationGenerator) generateDownMigration(tableNames []string) []string {
	var downStatements []string
	for i := len(tableNames) - 1; i >= 0; i-- {
		downStatements = append(downStatements, fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableNames[i]))
	}
	return downStatements
}

func (mg *MigrationGenerator) removeDuplicateTableDefinitions(sqlString string) string {
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

func (mg *MigrationGenerator) generateCreateTableSQL(stmt *germ.Statement) (string, []string, error) {
	fields := mg.extractFields(stmt)

	var fieldDefs, constraints []string
	var dependencies []string
	for _, field := range fields {
		if strings.Contains(field, "CONSTRAINT") {
			constraints = append(constraints, field)
			// Extract dependency from foreign key constraint
			re := regexp.MustCompile(`REFERENCES (\w+)`)
			matches := re.FindStringSubmatch(field)
			if len(matches) > 1 {
				dependencies = append(dependencies, matches[1])
			}
		} else {
			fieldDefs = append(fieldDefs, field)
		}
	}

	fieldDefsSQL := strings.Join(fieldDefs, ",\n ")
	constraintsSQL := strings.Join(constraints, ",\n ")

	var createTableSQL string
	if constraintsSQL != "" {
		createTableSQL = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s,\n %s\n);",
			stmt.Table, fieldDefsSQL, constraintsSQL)
	} else {
		createTableSQL = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s\n);",
			stmt.Table, fieldDefsSQL)
	}

	return createTableSQL, dependencies, nil
}

func (mg *MigrationGenerator) extractFields(stmt *germ.Statement) []string {
	var fields []string

	for _, field := range stmt.Schema.Fields {
		fieldSQL := mg.db.Migrator().FullDataTypeOf(field).SQL
		if fieldSQL != "" {
			fieldStr := mg.generateFieldDefinition(field, fieldSQL)
			fields = append(fields, fieldStr)

			if fk, ok := mg.generateForeignKeyConstraint(stmt, field); ok {
				fields = append(fields, fk)
			}

			if uni := mg.generateUniqueConstraint(stmt, field); uni != "" {
				fields = append(fields, uni)
			}
		}
	}

	return fields
}

func (mg *MigrationGenerator) generateFieldDefinition(field *schema.Field, fieldSQL string) string {
	fieldStr := fmt.Sprintf("%s %s", field.DBName, fieldSQL)

	if field.AutoIncrement {
		if mg.db.Name() == "sqlite" || mg.db.Name() == "sqlite3" {
			fieldStr = fmt.Sprintf("%s INTEGER PRIMARY KEY AUTOINCREMENT", field.DBName)
		} else if mg.db.Name() == "mysql" {
			fieldStr = fmt.Sprintf("%s INT PRIMARY KEY AUTO_INCREMENT", field.DBName)
		}
	}

	if field.NotNull {
		fieldStr = mg.removeNotNull(fieldStr)
		fieldStr += " NOT NULL"
	}

	// Ensure DATETIME for time.Time fields in SQLite
	if strings.Contains(fieldSQL, "datetime") {
		if mg.db.Name() == "sqlite" || mg.db.Name() == "sqlite3" {
			fieldStr = fmt.Sprintf("%s DATETIME", field.DBName) // Set to DATETIME instead of TEXT
		}
	} else if mg.db.Name() == "mysql" && strings.Contains(fieldSQL, "TEXT") {
		fieldStr = fmt.Sprintf("%s VARCHAR(255)", field.DBName)
	}

	return fieldStr
}

func (mg *MigrationGenerator) generateForeignKeyConstraint(stmt *germ.Statement, field *schema.Field) (string, bool) {
	if field.TagSettings["FOREIGNKEY"] != "" {
		refTable, refField := mg.parseTagSetting(field.TagSettings["FOREIGNKEY"])
		if refTable == "" {
			refTable = strings.TrimSuffix(field.DBName, "ID")
			if refTable == "" {
				refTable = field.DBName
			}
		}
		if refField == "" {
			refField = "ID"
		}

		fkName := fmt.Sprintf("fk_%s_%s", stmt.Table, field.DBName)
		return fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE CASCADE",
			fkName, field.DBName, refTable, refField), true
	}
	return "", false
}

func (mg *MigrationGenerator) generateUniqueConstraint(stmt *germ.Statement, field *schema.Field) string {
	if field.Unique {
		uniName := fmt.Sprintf("uk_%s_%s", stmt.Table, field.DBName)
		return fmt.Sprintf("CONSTRAINT %s UNIQUE (%s)", uniName, field.DBName)
	}
	return ""
}

func (mg *MigrationGenerator) parseTagSetting(tagSetting string) (string, string) {
	parts := strings.Split(tagSetting, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

func (mg *MigrationGenerator) removeNotNull(fieldStr string) string {
	return strings.ReplaceAll(fieldStr, " NOT NULL", "")
}

func (mg *MigrationGenerator) generateJoinTableSQL(stmt *germ.Statement) (string, []string, error) {
	var joinTableSQL []string
	var dependencies []string

	for _, rel := range stmt.Schema.Relationships.Relations {
		if rel.Type == schema.Many2Many && rel.JoinTable != nil {
			joinTable := rel.JoinTable
			refTable := rel.FieldSchema.Table

			joinFields := []string{
				fmt.Sprintf("%s_id INTEGER NOT NULL", stmt.Table),
				fmt.Sprintf("%s_id INTEGER NOT NULL", rel.FieldSchema.Table),
				fmt.Sprintf("PRIMARY KEY (%s_id, %s_id)", stmt.Table, rel.FieldSchema.Table),
				fmt.Sprintf("FOREIGN KEY (%s_id) REFERENCES %s(id) ON DELETE CASCADE", stmt.Table, refTable),
				fmt.Sprintf("FOREIGN KEY (%s_id) REFERENCES %s(id) ON DELETE CASCADE", rel.FieldSchema.Table, refTable),
			}

			joinTableSQL = append(joinTableSQL, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n %s\n);",
				joinTable.Name, strings.Join(joinFields, ",\n ")))

			dependencies = append(dependencies, stmt.Table, rel.FieldSchema.Table)
		}
	}

	return strings.Join(joinTableSQL, "\n\n"), dependencies, nil
}
