// admin/urls.go
package admin

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/admin", (&UserView{}).Index)
}
