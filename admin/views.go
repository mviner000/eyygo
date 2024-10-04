// admin/views.go
package admin

import (
	"github.com/gofiber/fiber/v2"
)

type UserView struct{}

func (u *UserView) Index(c *fiber.Ctx) error {
	users := GetAllUsers()
	return c.Render("admin/templates/users", fiber.Map{
		"users": users,
	})
}
