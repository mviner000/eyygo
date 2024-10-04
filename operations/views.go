// views.go
package operations

import (
	"fmt"
)

// CreateView represents an operation to create a view in the database.
type CreateView struct {
	ViewName string
	Query    string
}

// Execute performs the operation to create a view with the specified query.
func (c *CreateView) Execute() error {
	fmt.Printf("Creating view %s\n", c.ViewName)
	sql := fmt.Sprintf("CREATE VIEW %s AS %s;", c.ViewName, c.Query)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropView represents an operation to drop a view from the database.
type DropView struct {
	ViewName string
}

// Execute performs the operation to drop the specified view.
func (d *DropView) Execute() error {
	fmt.Printf("Dropping view %s\n", d.ViewName)
	sql := fmt.Sprintf("DROP VIEW %s;", d.ViewName)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// UpdateView represents an operation to update a view in the database.
type UpdateView struct {
	ViewName string
	NewQuery string
}

// Execute performs the operation to update the specified view with a new query.
func (u *UpdateView) Execute() error {
	fmt.Printf("Updating view %s\n", u.ViewName)
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s;", u.ViewName, u.NewQuery)
	fmt.Println("Executing SQL:", sql)
	return nil
}
