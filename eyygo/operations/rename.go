package operations

import (
	"fmt"
)

type RenameModel struct {
	OldName string
	NewName string
}

func (r *RenameModel) Execute() error {
	fmt.Printf("Renaming model from %s to %s\n", r.OldName, r.NewName)
	sql := fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", r.OldName, r.NewName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

type RenameField struct {
	ModelName    string
	OldFieldName string
	NewFieldName string
}

func (r *RenameField) Execute() error {
	fmt.Printf("Renaming field from %s to %s in model %s\n", r.OldFieldName, r.NewFieldName, r.ModelName)
	sql := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;", r.ModelName, r.OldFieldName, r.NewFieldName)
	fmt.Println("Executing SQL:", sql)
	return nil
}
