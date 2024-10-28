// admin/registry.go
package admin

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// ModelAdmin defines the configuration and customization for a model in the admin interface
type ModelAdmin struct {
	Model        interface{}
	ListFields   []string
	SearchFields []string
	FilterFields []string
	OrderFields  []string
	FormFields   []string
	DB           *gorm.DB
}

// AdminSite handles the registration and management of models
type AdminSite struct {
	registry map[string]*ModelAdmin
	db       *gorm.DB
}

// NewAdminSite creates a new AdminSite instance
func NewAdminSite(db *gorm.DB) *AdminSite {
	return &AdminSite{
		registry: make(map[string]*ModelAdmin),
		db:       db,
	}
}

// Register adds a model to the admin interface
func (site *AdminSite) Register(model interface{}, config *ModelAdmin) {
	modelType := reflect.TypeOf(model)
	modelName := strings.ToLower(modelType.Name())

	if config == nil {
		config = &ModelAdmin{
			Model: model,
			DB:    site.db,
		}

		// Auto-generate fields if not specified
		val := reflect.ValueOf(model)
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.Name != "Model" && !strings.HasSuffix(field.Name, "At") {
				config.ListFields = append(config.ListFields, field.Name)
				config.FormFields = append(config.FormFields, field.Name)
			}
		}
	}

	site.registry[modelName] = config
}

// GetModelAdmin retrieves the ModelAdmin configuration for a given model name
func (site *AdminSite) GetModelAdmin(modelName string) (*ModelAdmin, bool) {
	admin, exists := site.registry[modelName]
	return admin, exists
}

// GetRegisteredModels returns all registered models
func (site *AdminSite) GetRegisteredModels() map[string]*ModelAdmin {
	return site.registry
}
