// handlers/admin.go
package handlers

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyygo/admin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// ListModels returns all registered models in the admin interface
func (h *AdminHandler) ListModels(c *fiber.Ctx) error {
	models := admin.Site.GetRegisteredModels()
	return c.JSON(fiber.Map{
		"models": models,
	})
}

// ListModelEntries returns paginated entries for a specific model
func (h *AdminHandler) ListModelEntries(c *fiber.Ctx) error {
	modelName := c.Params("model")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)

	modelAdmin, exists := admin.Site.GetModelAdmin(modelName)
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	var results []interface{}
	var count int64

	query := modelAdmin.DB.Model(modelAdmin.Model)
	query.Count(&count)

	offset := (page - 1) * perPage
	err := query.Offset(offset).Limit(perPage).Find(&results).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch entries",
		})
	}

	return c.JSON(fiber.Map{
		"data":        results,
		"total":       count,
		"page":        page,
		"total_pages": (count + int64(perPage) - 1) / int64(perPage),
	})
}

// GetModelEntry returns a specific model entry
func (h *AdminHandler) GetModelEntry(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")

	modelAdmin, exists := admin.Site.GetModelAdmin(modelName)
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	var result interface{}
	err := modelAdmin.DB.First(&result, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Entry not found",
		})
	}

	return c.JSON(result)
}

// CreateModelEntry creates a new model entry
func (h *AdminHandler) CreateModelEntry(c *fiber.Ctx) error {
	modelName := c.Params("model")

	modelAdmin, exists := admin.Site.GetModelAdmin(modelName)
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	entry := reflect.New(reflect.TypeOf(modelAdmin.Model).Elem()).Interface()
	if err := c.BodyParser(entry); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := modelAdmin.DB.Create(entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create entry",
		})
	}

	return c.JSON(entry)
}

// UpdateModelEntry updates a specific model entry
func (h *AdminHandler) UpdateModelEntry(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")

	modelAdmin, exists := admin.Site.GetModelAdmin(modelName)
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	entry := reflect.New(reflect.TypeOf(modelAdmin.Model).Elem()).Interface()
	if err := modelAdmin.DB.First(entry, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Entry not found",
		})
	}

	if err := c.BodyParser(entry); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := modelAdmin.DB.Save(entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update entry",
		})
	}

	return c.JSON(entry)
}

// DeleteModelEntry deletes a specific model entry
func (h *AdminHandler) DeleteModelEntry(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")

	modelAdmin, exists := admin.Site.GetModelAdmin(modelName)
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	entry := reflect.New(reflect.TypeOf(modelAdmin.Model).Elem()).Interface()
	if err := modelAdmin.DB.Delete(entry, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete entry",
		})
	}

	return c.SendStatus(204)
}
