package admin

import (
	"github.com/gofiber/fiber/v2"
	customhttp "github.com/mviner000/eyymi/http"
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
	if err := v.DB.Find(&users).Error; err != nil {
		return customhttp.JsonResponse(fiber.Map{"error": "Error fetching users"}, fiber.StatusInternalServerError, nil).Render(c)
	}
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
		return customhttp.JsonResponse(fiber.Map{"error": "Invalid input"}, fiber.StatusBadRequest, nil).Render(c)
	}
	if err := v.DB.Create(user).Error; err != nil {
		return customhttp.JsonResponse(fiber.Map{"error": "Error creating user"}, fiber.StatusInternalServerError, nil).Render(c)
	}
	return v.UserList(c) // Return the updated user list
}

func (v *AdminViews) UserEdit(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := v.DB.First(&user, id).Error; err != nil {
		return customhttp.JsonResponse(fiber.Map{"error": "User not found"}, fiber.StatusNotFound, nil).Render(c)
	}
	return c.Render("admin/templates/user_form", fiber.Map{
		"Title": "Edit User",
		"User":  user,
	})
}

func (v *AdminViews) UserUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return customhttp.JsonResponse(fiber.Map{"error": "Invalid input"}, fiber.StatusBadRequest, nil).Render(c)
	}
	if err := v.DB.Model(&User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return customhttp.JsonResponse(fiber.Map{"error": "Error updating user"}, fiber.StatusInternalServerError, nil).Render(c)
	}
	return v.UserList(c) // Return the updated user list
}

func (v *AdminViews) UserDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := v.DB.Delete(&User{}, id).Error; err != nil {
		return customhttp.JsonResponse(fiber.Map{"error": "Error deleting user"}, fiber.StatusInternalServerError, nil).Render(c)
	}
	// Return an empty JSON object instead of no content
	return customhttp.JsonResponse(fiber.Map{}, fiber.StatusOK, nil).Render(c)
}
