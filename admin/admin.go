package admin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewAdminModule(db *gorm.DB) (*AdminViews, func(*fiber.App)) {
	views := NewAdminViews(db)
	return views, func(app *fiber.App) {
		SetupRoutes(app, views)
	}
}
