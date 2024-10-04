// index.go
package operations

import (
	"fmt"
	"strings"
)

// CreateIndex represents an operation to create an index on a database table.
type CreateIndex struct {
	TableName string
	IndexName string
	Columns   []string
}

// Execute performs the operation to create an index on the specified table.
func (c *CreateIndex) Execute() error {
	fmt.Printf("Creating index %s on table %s\n", c.IndexName, c.TableName)
	sql := fmt.Sprintf("CREATE INDEX %s ON %s (%s);", c.IndexName, c.TableName, strings.Join(c.Columns, ", "))
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropIndex represents an operation to drop an index from a database table.
type DropIndex struct {
	TableName string
	IndexName string
}

// Execute performs the operation to drop the index from the specified table.
func (d *DropIndex) Execute() error {
	fmt.Printf("Dropping index %s from table %s\n", d.IndexName, d.TableName)
	sql := fmt.Sprintf("DROP INDEX %s ON %s;", d.IndexName, d.TableName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// RenameIndex represents an operation to rename an index on a database table.
type RenameIndex struct {
	TableName    string
	OldIndexName string
	NewIndexName string
}

// Execute performs the operation to rename the index on the specified table.
func (r *RenameIndex) Execute() error {
	fmt.Printf("Renaming index %s to %s on table %s\n", r.OldIndexName, r.NewIndexName, r.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s RENAME INDEX %s TO %s;", r.TableName, r.OldIndexName, r.NewIndexName)
	fmt.Println("Executing SQL:", sql)
	return nil
}
