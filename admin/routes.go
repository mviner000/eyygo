package admin

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, views *AdminViews) {
	admin := app.Group("/admin")
	admin.Get("/", views.Dashboard)
	admin.Get("/users", views.UserList)
	admin.Get("/users/new", views.UserCreate)
	admin.Post("/users", views.UserStore)
	admin.Get("/users/:id", views.UserEdit)
	admin.Put("/users/:id", views.UserUpdate)
	admin.Delete("/users/:id", views.UserDelete)
}
