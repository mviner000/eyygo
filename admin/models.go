// admin/models.go
package admin

import (
	"github.com/mviner000/eyymi/fields"
	"github.com/mviner000/eyymi/operations"
)

type User struct {
	Username string
}

func main() {
	// Create a new model for the User and add a CharField in one line
	userModel := operations.NewModel("users")
	userModel.AddField("username", fields.CharField("username", 255))
}

func GetAllUsers() []User {
	// This is a stub. Replace it with actual database fetching logic.
	return []User{
		{Username: "alice"},
		{Username: "bob"},
	}
}
