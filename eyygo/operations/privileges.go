// privileges.go
package operations

import (
	"fmt"
	"strings"
)

// GrantPrivileges represents an operation to grant privileges to a user.
type GrantPrivileges struct {
	User       string
	Privileges []string
	On         string
}

// Execute performs the grant operation for the specified user and privileges.
func (g *GrantPrivileges) Execute() error {
	fmt.Printf("Granting privileges to user %s\n", g.User)
	sql := fmt.Sprintf("GRANT %s ON %s TO '%s';", strings.Join(g.Privileges, ", "), g.On, g.User)
	fmt.Println("Executing SQL:", sql)
	return nil
}

// RevokePrivileges represents an operation to revoke privileges from a user.
type RevokePrivileges struct {
	User       string
	Privileges []string
	On         string
}

// Execute performs the revoke operation for the specified user and privileges.
func (r *RevokePrivileges) Execute() error {
	fmt.Printf("Revoking privileges from user %s\n", r.User)
	sql := fmt.Sprintf("REVOKE %s ON %s FROM '%s';", strings.Join(r.Privileges, ", "), r.On, r.User)
	fmt.Println("Executing SQL:", sql)
	return nil
}
