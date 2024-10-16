// columns.go
package operations

import (
	"fmt"
)

// AddColumn represents an operation to add a column to a table.
type AddColumn struct {
	TableName  string
	ColumnName string
	ColumnType string
}

// Execute performs the operation to add a column to the specified table.
func (a *AddColumn) Execute() error {
	fmt.Printf("Adding column %s to table %s\n", a.ColumnName, a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", a.TableName, a.ColumnName, a.ColumnType)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropColumn represents an operation to drop a column from a table.
type DropColumn struct {
	TableName  string
	ColumnName string
}

// Execute performs the operation to drop the column from the specified table.
func (d *DropColumn) Execute() error {
	fmt.Printf("Dropping column %s from table %s\n", d.ColumnName, d.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", d.TableName, d.ColumnName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// RenameColumn represents an operation to rename a column in a table.
type RenameColumn struct {
	TableName     string
	OldColumnName string
	NewColumnName string
}

// Execute performs the operation to rename the column in the specified table.
func (r *RenameColumn) Execute() error {
	fmt.Printf("Renaming column %s to %s in table %s\n", r.OldColumnName, r.NewColumnName, r.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;", r.TableName, r.OldColumnName, r.NewColumnName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// ModifyColumnType represents an operation to change the type of a column in a table.
type ModifyColumnType struct {
	TableName  string
	ColumnName string
	NewType    string
}

// Execute performs the operation to change the column type in the specified table.
func (m *ModifyColumnType) Execute() error {
	fmt.Printf("Modifying column %s type to %s in table %s\n", m.ColumnName, m.NewType, m.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;", m.TableName, m.ColumnName, m.NewType)
	fmt.Println("Executing SQL:", sql)
	return nil
}
