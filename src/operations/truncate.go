// truncate.go
package operations

import (
	"fmt"
)

// TruncateTable represents an operation to truncate a database table.
type TruncateTable struct {
	TableName string
}

// Execute performs the truncate operation on the specified table.
func (t *TruncateTable) Execute() error {
	fmt.Printf("Truncating table %s\n", t.TableName)
	sql := fmt.Sprintf("TRUNCATE TABLE %s;", t.TableName)
	fmt.Println("Executing SQL:", sql)
	return nil
}
