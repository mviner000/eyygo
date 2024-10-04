package operations

import (
	"fmt"
)

type RemoveConstraint struct {
	ModelName      string
	ConstraintName string
}

func (r *RemoveConstraint) Execute() error {
	fmt.Printf("Removing constraint %s from model %s\n", r.ConstraintName, r.ModelName)
	sql := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;", r.ModelName, r.ConstraintName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

type RemoveIndex struct {
	ModelName string
	IndexName string
}

func (r *RemoveIndex) Execute() error {
	fmt.Printf("Removing index %s from model %s\n", r.IndexName, r.ModelName)
	sql := fmt.Sprintf("DROP INDEX %s ON %s;", r.IndexName, r.ModelName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

type RemoveField struct {
	ModelName string
	FieldName string
}

func (r *RemoveField) Execute() error {
	fmt.Printf("Removing field %s from model %s\n", r.FieldName, r.ModelName)
	sql := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", r.ModelName, r.FieldName)
	fmt.Println("Executing SQL:", sql)
	return nil
}
