// admin/routes.go
package admin

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	userView := &UserView{}

	app.Get("/admin/users", userView.Index)
	app.Get("/admin/users/create", userView.Create)
	app.Post("/admin/users", userView.Store)
}
