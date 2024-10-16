package conf

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/project_name"
	"github.com/mviner000/eyymi/src/registry"
)

var Settings *project_name.SettingsStruct
var modelsRegistered bool

func init() {
	LoadSettings()
}

// LoadSettings reloads the settings from the project_name package
func LoadSettings() {
	project_name.LoadSettings()
	Settings = &project_name.AppSettings
}

// GetSettings returns the global settings
func GetSettings() *project_name.SettingsStruct {
	return Settings
}

// GetFullProjectName returns the FullProjectName directly
func GetFullProjectName() string {
	return Settings.FullProjectName
}

// EnsureModelsRegistered ensures that models are registered
func EnsureModelsRegistered() {
	if !modelsRegistered {
		project_name.RegisterModels()
		modelsRegistered = true
	}
}

// GetRegisteredModelNames returns the names of all registered models
func GetRegisteredModelNames() []string {
	EnsureModelsRegistered()
	return registry.GetRegisteredModelNames()
}

// GetRegisteredModel returns a registered model by name
func GetRegisteredModel(name string) (interface{}, bool) {
	EnsureModelsRegistered()
	return registry.GetRegisteredModel(name)
}

// GetRegisteredModelsInfo returns a string with information about all registered models
func GetRegisteredModelsInfo() string {
	EnsureModelsRegistered()
	return registry.GetRegisteredModels()
}

// App interface definition
type App interface {
	SetupRoutes(app *fiber.App)
}

// ProjectNameApp struct
type ProjectNameApp struct {
	settings *project_name.SettingsStruct
}

// SetupRoutes implementation for ProjectNameApp
func (a *ProjectNameApp) SetupRoutes(app *fiber.App) {
	// Setup routes using a.settings
	// This is where you'd implement the actual route setup for project_name
}

// NewApp creates and returns a new App instance
func NewApp() App {
	return &ProjectNameApp{settings: Settings}
}
