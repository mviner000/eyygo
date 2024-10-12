package admin

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/eyygo/auth"
)

// SetupRoutes sets up the admin routes
func SetupRoutes(app *fiber.App) {
	log.Println("Admin: Starting to set up admin routes")
	adminGroup := app.Group("/admin")

	// Public routes
	log.Println("Admin: Setting up public routes")
	adminGroup.Get("/login", func(c *fiber.Ctx) error {
		log.Println("Admin: LoginForm handler called")

		return LoginForm(c)
	})
	adminGroup.Post("/login", func(c *fiber.Ctx) error {
		log.Println("Admin: Login handler called")
		return Login(c)
	})

	// Protected routes
	log.Println("Admin: Setting up protected routes with AuthMiddleware")
	adminGroup.Get("/dashboard", logMiddleware("Dashboard"), auth.AuthMiddleware, Dashboard)
	adminGroup.Get("/users", logMiddleware("UserIndex"), auth.AuthMiddleware, UserList)
	adminGroup.Get("/users/create", logMiddleware("UserCreate"), auth.AuthMiddleware, UserCreate)
	adminGroup.Post("/users", logMiddleware("UserStore"), auth.AuthMiddleware, UserStore)
	adminGroup.Post("/logout", logMiddleware("Logout"), auth.AuthMiddleware, Logout)

	log.Println("Admin: Finished setting up admin routes")
}

// logMiddleware is a helper function to log when a route is accessed
func logMiddleware(routeName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("Admin: Accessing %s route", routeName)
		return c.Next()
	}
}
