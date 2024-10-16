package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/src/auth"
)

func UserList(c *fiber.Ctx) error {
	users, err := auth.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving user list")
	}

	return c.Render("src/admin/templates/user_list", fiber.Map{
		"Users": users,
	})
}
