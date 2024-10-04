// alter.go
package operations

import (
	"fmt"
)

// AlterTableComment represents an operation to change the comment of a table.
type AlterTableComment struct {
	TableName string
	Comment   string
}

// Execute performs the operation to change the comment of the specified table.
func (a *AlterTableComment) Execute() error {
	fmt.Printf("Changing comment for table %s\n", a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s COMMENT = '%s';", a.TableName, a.Comment)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// AlterTableEngine represents an operation to change the storage engine of a table.
type AlterTableEngine struct {
	TableName string
	Engine    string
}

// Execute performs the operation to change the storage engine of the specified table.
func (a *AlterTableEngine) Execute() error {
	fmt.Printf("Changing storage engine for table %s\n", a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s ENGINE = %s;", a.TableName, a.Engine)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// AlterTableCollation represents an operation to change the collation of a table.
type AlterTableCollation struct {
	TableName string
	Collation string
}

// Execute performs the operation to change the collation of the specified table.
func (a *AlterTableCollation) Execute() error {
	fmt.Printf("Changing collation for table %s\n", a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s COLLATE = %s;", a.TableName, a.Collation)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// AlterTableCharacterSet represents an operation to change the character set of a table.
type AlterTableCharacterSet struct {
	TableName    string
	CharacterSet string
}

// Execute performs the operation to change the character set of the specified table.
func (a *AlterTableCharacterSet) Execute() error {
	fmt.Printf("Changing character set for table %s\n", a.TableName)
	sql := fmt.Sprintf("ALTER TABLE %s CHARACTER SET = %s;", a.TableName, a.CharacterSet)
	fmt.Println("Executing SQL:", sql)
	return nil
}
