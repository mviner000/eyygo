// triggers.go
package operations

import (
	"fmt"
)

// CreateTrigger represents an operation to create a trigger in the database.
type CreateTrigger struct {
	TriggerName string
	Timing      string
	Event       string
	TableName   string
	Statement   string
}

// Execute performs the operation to create a trigger with the specified properties.
func (c *CreateTrigger) Execute() error {
	fmt.Printf("Creating trigger %s\n", c.TriggerName)
	sql := fmt.Sprintf("CREATE TRIGGER %s %s %s ON %s FOR EACH ROW %s;",
		c.TriggerName, c.Timing, c.Event, c.TableName, c.Statement)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// DropTrigger represents an operation to drop a trigger from the database.
type DropTrigger struct {
	TriggerName string
}

// Execute performs the operation to drop the specified trigger.
func (d *DropTrigger) Execute() error {
	fmt.Printf("Dropping trigger %s\n", d.TriggerName)
	sql := fmt.Sprintf("DROP TRIGGER %s;", d.TriggerName)
	fmt.Println("Executing SQL:", sql)
	return nil
}
