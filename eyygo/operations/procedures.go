// procedures.go
package operations

import (
	"fmt"
)

// CreateProcedure represents an operation to create a stored procedure in the database.
type CreateProcedure struct {
	ProcedureName string
	Parameters    string
	Body          string
}

// Execute performs the operation to create a stored procedure with the specified properties.
func (c *CreateProcedure) Execute() error {
	fmt.Printf("Creating procedure %s\n", c.ProcedureName)
	sql := fmt.Sprintf("CREATE PROCEDURE %s(%s) AS $$ BEGIN %s END; $$ LANGUAGE plpgsql;",
		c.ProcedureName, c.Parameters, c.Body)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropProcedure represents an operation to drop a stored procedure from the database.
type DropProcedure struct {
	ProcedureName string
	Parameters    string
}

// Execute performs the operation to drop the specified stored procedure.
func (d *DropProcedure) Execute() error {
	fmt.Printf("Dropping procedure %s\n", d.ProcedureName)
	sql := fmt.Sprintf("DROP PROCEDURE %s(%s);", d.ProcedureName, d.Parameters)
	fmt.Println("Executing SQL:", sql)
	return nil
}
