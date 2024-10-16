// tables.go
package operations

import (
	"fmt"
)

// CreateTable represents an operation to create a table in the database.
type CreateTable struct {
	TableName string
	Columns   []string
}

// Execute performs the operation to create a table with the specified columns.
func (c *CreateTable) Execute() error {
	fmt.Printf("Creating table %s\n", c.TableName)
	sql := fmt.Sprintf("CREATE TABLE %s (%s);", c.TableName, joinColumns(c.Columns))
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropTable represents an operation to drop a table from the database.
type DropTable struct {
	TableName string
}

// Execute performs the operation to drop the specified table.
func (d *DropTable) Execute() error {
	fmt.Printf("Dropping table %s\n", d.TableName)
	sql := fmt.Sprintf("DROP TABLE %s;", d.TableName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// RenameTable represents an operation to rename a table in the database.
type RenameTable struct {
	OldTableName string
	NewTableName string
}

// Execute performs the operation to rename the specified table.
func (r *RenameTable) Execute() error {
	fmt.Printf("Renaming table %s to %s\n", r.OldTableName, r.NewTableName)
	sql := fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", r.OldTableName, r.NewTableName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// Helper function to join columns into a single string.
func joinColumns(columns []string) string {
	return fmt.Sprintf("%s", columns)
}
