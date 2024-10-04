// functions.go
package operations

import (
	"fmt"
)

// CreateFunction represents an operation to create a function in the database.
type CreateFunction struct {
	FunctionName string
	Parameters   string
	ReturnType   string
	Body         string
}

// Execute performs the operation to create a function with the specified properties.
func (c *CreateFunction) Execute() error {
	fmt.Printf("Creating function %s\n", c.FunctionName)
	sql := fmt.Sprintf("CREATE FUNCTION %s(%s) RETURNS %s AS $$ BEGIN %s END; $$ LANGUAGE plpgsql;",
		c.FunctionName, c.Parameters, c.ReturnType, c.Body)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropFunction represents an operation to drop a function from the database.
type DropFunction struct {
	FunctionName string
	Parameters   string
}

// Execute performs the operation to drop the specified function.
func (d *DropFunction) Execute() error {
	fmt.Printf("Dropping function %s\n", d.FunctionName)
	sql := fmt.Sprintf("DROP FUNCTION %s(%s);", d.FunctionName, d.Parameters)
	fmt.Println("Executing SQL:", sql)
	return nil
}
