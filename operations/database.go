// database.go
package operations

import (
	"fmt"
)

// CreateDatabase represents an operation to create a new database.
type CreateDatabase struct {
	DatabaseName string
}

// Execute performs the operation to create a new database.
func (c *CreateDatabase) Execute() error {
	fmt.Printf("Creating database %s\n", c.DatabaseName)
	sql := fmt.Sprintf("CREATE DATABASE %s;", c.DatabaseName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropDatabase represents an operation to drop an existing database.
type DropDatabase struct {
	DatabaseName string
}

// Execute performs the operation to drop the specified database.
func (d *DropDatabase) Execute() error {
	fmt.Printf("Dropping database %s\n", d.DatabaseName)
	sql := fmt.Sprintf("DROP DATABASE %s;", d.DatabaseName)
	fmt.Println("Executing SQL:", sql)
	return nil
}
