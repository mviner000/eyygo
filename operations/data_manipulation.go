// data_manipulation.go
package operations

import (
	"fmt"
	"strings"
)

// UpdateRows represents an operation to update rows in a database table.
type UpdateRows struct {
	TableName   string
	SetClause   string
	WhereClause string
}

// Execute performs the update operation on the specified table.
func (u *UpdateRows) Execute() error {
	fmt.Printf("Updating rows in table %s\n", u.TableName)
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s;", u.TableName, u.SetClause, u.WhereClause)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// InsertRows represents an operation to insert rows into a database table.
type InsertRows struct {
	TableName string
	Columns   []string
	Values    []string
}

// Execute performs the insert operation on the specified table.
func (i *InsertRows) Execute() error {
	fmt.Printf("Inserting rows into table %s\n", i.TableName)
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", i.TableName, strings.Join(i.Columns, ", "), strings.Join(i.Values, ", "))
	fmt.Println("Executing SQL:", sql)
	return nil
}

// SelectRows represents an operation to select rows from a database table.
type SelectRows struct {
	TableName   string
	Columns     []string
	WhereClause string
}

// Execute performs the select operation on the specified table.
func (s *SelectRows) Execute() error {
	fmt.Printf("Selecting rows from table %s\n", s.TableName)
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s;", strings.Join(s.Columns, ", "), s.TableName, s.WhereClause)
	fmt.Println("Executing SQL:", sql)
	return nil
}
