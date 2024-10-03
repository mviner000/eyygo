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
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", a.ModelName, a.FieldName, a.FieldType)
	fmt.Println("Executing SQL:", sql)
	return nil
}

type AddIndex struct {
	ModelName  string
	IndexName  string
	FieldNames []string
}

func (a *AddIndex) Execute() error {
	fmt.Printf("Adding index %s to model %s\n", a.IndexName, a.ModelName)
	sql := fmt.Sprintf("CREATE INDEX %s ON %s (%s);", a.IndexName, a.ModelName, strings.Join(a.FieldNames, ", "))
	fmt.Println("Executing SQL:", sql)
	return nil
}

type AddConstraint struct {
	ModelName      string
	ConstraintName string
	ConstraintSQL  string
}

func (a *AddConstraint) Execute() error {
	fmt.Printf("Adding constraint %s to model %s\n", a.ConstraintName, a.ModelName)
	sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s;", a.ModelName, a.ConstraintName, a.ConstraintSQL)
	fmt.Println("Executing SQL:", sql)
	return nil
}
