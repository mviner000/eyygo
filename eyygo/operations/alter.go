package operations

import (
	"fmt"
	"strings"
)

// AlterModelTable represents an operation to change the name of a database table
type AlterModelTable struct {
	Name    string
	NewName string
}

func (a *AlterModelTable) Execute() error {
	fmt.Printf("Changing table name for model %s to %s\n", a.Name, a.NewName)
	sql := fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", a.Name, a.NewName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// AlterUniqueTogether represents an operation to change the unique_together constraints
type AlterUniqueTogether struct {
	ModelName string
	Fields    [][]string
}

func (a *AlterUniqueTogether) Execute() error {
	fmt.Printf("Altering unique together constraints for model %s\n", a.ModelName)
	for _, constraint := range a.Fields {
		sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT unique_%s UNIQUE (%s);",
			a.ModelName, strings.Join(constraint, "_"), strings.Join(constraint, ", "))
		fmt.Println("Executing SQL:", sql)
	}
	return nil
}

// AlterIndexTogether represents an operation to change the index_together constraints
type AlterIndexTogether struct {
	ModelName string
	Fields    [][]string
}

func (a *AlterIndexTogether) Execute() error {
	fmt.Printf("Altering index together for model %s\n", a.ModelName)
	for i, indexFields := range a.Fields {
		sql := fmt.Sprintf("CREATE INDEX idx_%s_%d ON %s (%s);",
			a.ModelName, i, a.ModelName, strings.Join(indexFields, ", "))
		fmt.Println("Executing SQL:", sql)
	}
	return nil
}

// AlterOrderWithRespectTo represents an operation to change the order_with_respect_to option
type AlterOrderWithRespectTo struct {
	ModelName string
	Field     string
}

func (a *AlterOrderWithRespectTo) Execute() error {
	fmt.Printf("Altering order with respect to %s for model %s\n", a.Field, a.ModelName)
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN _order INTEGER;", a.ModelName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// AlterModelOptions represents an operation to change the options on a model
type AlterModelOptions struct {
	ModelName string
	Options   map[string]interface{}
}

func (a *AlterModelOptions) Execute() error {
	fmt.Printf("Altering model options for %s\n", a.ModelName)
	for option, value := range a.Options {
		fmt.Printf("Setting %s to %v\n", option, value)
		// Note: This doesn't typically result in SQL changes,
		// it's more about changing the model's metadata in the ORM
	}
	return nil
}

// AlterField represents an operation to change the definition of a field on a model
type AlterField struct {
	ModelName string
	FieldName string
	NewType   string
}

func (a *AlterField) Execute() error {
	fmt.Printf("Altering field %s on model %s\n", a.FieldName, a.ModelName)
	sql := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;",
		a.ModelName, a.FieldName, a.NewType)
	fmt.Println("Executing SQL:", sql)
	return nil
}
