package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/eyygo/http"
)

func UserCreate(c *fiber.Ctx) error {
	response := http.HttpResponseOK(fiber.Map{}, nil, "eyygo/admin/templates/user_form")
	return response.Render(c)
}

func UserStore(c *fiber.Ctx) error {
	return c.SendString("User creation logic not implemented")
}
