package admin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminViews struct {
	DB *gorm.DB
}

func NewAdminViews(db *gorm.DB) *AdminViews {
	return &AdminViews{DB: db}
}

func (v *AdminViews) Dashboard(c *fiber.Ctx) error {
	return c.Render("admin/templates/dashboard", fiber.Map{
		"Title": "Admin Dashboard",
	})
}

func (v *AdminViews) UserList(c *fiber.Ctx) error {
	var users []User
	v.DB.Find(&users)
	return c.Render("admin/templates/user_list", fiber.Map{
		"Title": "User List",
		"Users": users,
	})
}

func (v *AdminViews) UserCreate(c *fiber.Ctx) error {
	return c.Render("admin/templates/user_form", fiber.Map{
		"Title": "Create User",
	})
}

func (v *AdminViews) UserStore(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	v.DB.Create(user)
	return c.Redirect("/admin/users")
}

func (v *AdminViews) UserEdit(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	v.DB.First(&user, id)
	return c.Render("admin/templates/user_form", fiber.Map{
		"Title": "Edit User",
		"User":  user,
	})
}

func (v *AdminViews) UserUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	v.DB.Model(&User{}).Where("id = ?", id).Updates(user)
	return c.Redirect("/admin/users")
}

func (v *AdminViews) UserDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	v.DB.Delete(&User{}, id)
	return c.Redirect("/admin/users")
}
