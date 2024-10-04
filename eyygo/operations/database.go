// operations/database.go
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
	return nil
}

func (c *CreateDatabase) SQL() (string, error) {
	return fmt.Sprintf("CREATE DATABASE %s;", c.DatabaseName), nil
}

// DropDatabase represents an operation to drop an existing database.
type DropDatabase struct {
	DatabaseName string
}

// Execute performs the operation to drop the specified database.
func (d *DropDatabase) Execute() error {
	fmt.Printf("Dropping database %s\n", d.DatabaseName)
	return nil
}

func (d *DropDatabase) SQL() (string, error) {
	return fmt.Sprintf("DROP DATABASE %s;", d.DatabaseName), nil
}
