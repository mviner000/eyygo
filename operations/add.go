// operations/add.go
package operations

import (
	"fmt"
	"strings"
)

type AddField struct {
	ModelName string
	FieldName string
	FieldType string
}

func (a *AddField) Execute() error {
	fmt.Printf("Adding field %s to model %s\n", a.FieldName, a.ModelName)
	return nil
}

func (a *AddField) SQL() (string, error) {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", a.ModelName, a.FieldName, a.FieldType), nil
}

type AddIndex struct {
	ModelName  string
	IndexName  string
	FieldNames []string
}

func (a *AddIndex) Execute() error {
	fmt.Printf("Adding index %s to model %s\n", a.IndexName, a.ModelName)
	return nil
}

func (a *AddIndex) SQL() (string, error) {
	return fmt.Sprintf("CREATE INDEX %s ON %s (%s);", a.IndexName, a.ModelName, strings.Join(a.FieldNames, ", ")), nil
}

type AddConstraint struct {
	ModelName      string
	ConstraintName string
	ConstraintSQL  string
}

func (a *AddConstraint) Execute() error {
	fmt.Printf("Adding constraint %s to model %s\n", a.ConstraintName, a.ModelName)
	return nil
}

func (a *AddConstraint) SQL() (string, error) {
	return fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s;", a.ModelName, a.ConstraintName, a.ConstraintSQL), nil
}

type AddTable struct {
	Model *Model
}

func (a *AddTable) Execute() error {
	fmt.Printf("Adding table %s\n", a.Model.TableName)
	return nil
}

func (a *AddTable) SQL() (string, error) {
	fields := []string{}
	for _, field := range a.Model.Fields {
		fields = append(fields, fmt.Sprintf("%s %s", field.GetOptions().Name, field.SQLType()))
	}
	sql := fmt.Sprintf("CREATE TABLE %s (%s);", a.Model.TableName, strings.Join(fields, ", "))
	return sql, nil
}
