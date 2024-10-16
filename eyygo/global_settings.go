package conf

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/project_name"
)

var Settings *project_name.SettingsStruct

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
