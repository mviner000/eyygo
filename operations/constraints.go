// constraints.go
package operations

import (
	"fmt"
	"strings"
)

// AddForeignKeyConstraint represents an operation to add a foreign key constraint to a table.
type AddForeignKeyConstraint struct {
	TableName         string
	ConstraintName    string
	Columns           []string
	ReferencedTable   string
	ReferencedColumns []string
}

// Execute performs the operation to add a foreign key constraint to the specified table.
func (a *AddForeignKeyConstraint) Execute() error {
	fmt.Printf("Adding foreign key constraint %s to table %s\n", a.ConstraintName, a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s);",
		a.TableName, a.ConstraintName, strings.Join(a.Columns, ", "), a.ReferencedTable, strings.Join(a.ReferencedColumns, ", "))
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropForeignKeyConstraint represents an operation to drop a foreign key constraint from a table.
type DropForeignKeyConstraint struct {
	TableName      string
	ConstraintName string
}

// Execute performs the operation to drop the foreign key constraint from the specified table.
func (d *DropForeignKeyConstraint) Execute() error {
	fmt.Printf("Dropping foreign key constraint %s from table %s\n", d.ConstraintName, d.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s;", d.TableName, d.ConstraintName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// AddPrimaryKeyConstraint represents an operation to add a primary key constraint to a table.
type AddPrimaryKeyConstraint struct {
	TableName      string
	ConstraintName string
	Columns        []string
}

// Execute performs the operation to add a primary key constraint to the specified table.
func (a *AddPrimaryKeyConstraint) Execute() error {
	fmt.Printf("Adding primary key constraint %s to table %s\n", a.ConstraintName, a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY (%s);",
		a.TableName, a.ConstraintName, strings.Join(a.Columns, ", "))
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropPrimaryKeyConstraint represents an operation to drop a primary key constraint from a table.
type DropPrimaryKeyConstraint struct {
	TableName string
}

// Execute performs the operation to drop the primary key constraint from the specified table.
func (d *DropPrimaryKeyConstraint) Execute() error {
	fmt.Printf("Dropping primary key constraint from table %s\n", d.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY;", d.TableName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// AddUniqueConstraint represents an operation to add a unique constraint to a table.
type AddUniqueConstraint struct {
	TableName      string
	ConstraintName string
	Columns        []string
}

// Execute performs the operation to add a unique constraint to the specified table.
func (a *AddUniqueConstraint) Execute() error {
	fmt.Printf("Adding unique constraint %s to table %s\n", a.ConstraintName, a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s);",
		a.TableName, a.ConstraintName, strings.Join(a.Columns, ", "))
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropUniqueConstraint represents an operation to drop a unique constraint from a table.
type DropUniqueConstraint struct {
	TableName      string
	ConstraintName string
}

// Execute performs the operation to drop the unique constraint from the specified table.
func (d *DropUniqueConstraint) Execute() error {
	fmt.Printf("Dropping unique constraint %s from table %s\n", d.ConstraintName, d.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;", d.TableName, d.ConstraintName)
	fmt.Println("Executing SQL:", sql)
	return nil
}
