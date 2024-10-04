// admin/models.go
package admin

import (
	"fmt"

	"github.com/mviner000/eyymi/fields"
	"github.com/mviner000/eyymi/operations"
)

// GetModels defines and returns all the models for the application.
func GetModels() []*operations.Model {
	// Define a new model for the User
	userModel := operations.NewModel("users")
	userModel.AddField(fields.CharField("username", 255, fields.WithRequired(true)))
	userModel.AddField(fields.CharField("email", 255, fields.WithRequired(true), fields.WithUnique(true)))

	// Generate SQL for creating the table (for debugging purposes)
	sql := userModel.CreateTableSQL()
	fmt.Println(sql)

	// Return the models
	return []*operations.Model{userModel}
}

type User struct {
	Username string
	Email    string
}

// GetAllUsers returns a mock list of users
func GetAllUsers() []User {
	return []User{
		{Username: "user1", Email: "user1@example.com"},
		{Username: "user2", Email: "user2@example.com"},
	}
}
