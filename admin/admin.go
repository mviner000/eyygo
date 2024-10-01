package admin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) SetupRoutes(app *fiber.App) {
	admin := app.Group("/admin")
	admin.Get("/", h.Dashboard)
	admin.Get("/users", h.UserList)
	admin.Get("/users/new", h.UserCreate)
	admin.Post("/users", h.UserStore)
	admin.Get("/users/:id", h.UserEdit)
	admin.Put("/users/:id", h.UserUpdate)
	admin.Delete("/users/:id", h.UserDelete)
}

func (h *AdminHandler) Dashboard(c *fiber.Ctx) error {
	return c.Render("admin/dashboard", fiber.Map{
		"Title": "Admin Dashboard",
	})
}

func (h *AdminHandler) UserList(c *fiber.Ctx) error {
	var users []User
	h.DB.Find(&users)
	return c.Render("admin/user_list", fiber.Map{
		"Title": "User List",
		"Users": users,
	})
}

func (h *AdminHandler) UserCreate(c *fiber.Ctx) error {
	return c.Render("admin/user_form", fiber.Map{
		"Title": "Create User",
	})
}

func (h *AdminHandler) UserStore(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	h.DB.Create(user)
	return c.Redirect("/admin/users")
}

func (h *AdminHandler) UserEdit(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	h.DB.First(&user, id)
	return c.Render("admin/user_form", fiber.Map{
		"Title": "Edit User",
		"User":  user,
	})
}

func (h *AdminHandler) UserUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	h.DB.Model(&User{}).Where("id = ?", id).Updates(user)
	return c.Redirect("/admin/users")
}

func (h *AdminHandler) UserDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	h.DB.Delete(&User{}, id)
	return c.Redirect("/admin/users")
}
