package core

import (
	"fmt"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/eyygo/admin"
	"github.com/mviner000/eyymi/eyygo/monitor"
	"github.com/mviner000/eyymi/project_name"
)

var exampleapp App = nil

// Mocked INSTALLED_APPS-like structure
var INSTALLED_APPS = map[string]bool{
	"project_name": true,  // This app is enabled
	"otherapp":     false, // This app is disabled and won't be set up
	// Add more apps as needed
}

func getAppPackage(appName string) (App, error) {
	switch appName {
	case "project_name":
		return &project_name.AppName{}, nil
	case "exampleapp":
		if exampleapp != nil {
			return exampleapp, nil
		}
		return nil, fmt.Errorf("exampleapp is not available")
	default:
		return nil, fmt.Errorf("unknown app: %s", appName)
	}
}

func setupAppRoutes(app *fiber.App, appName string) {
	appPackage, err := getAppPackage(appName)
	if err != nil {
		if project_name.AppSettings.Debug {
			appLogger.Printf("Error setting up app: %v", err)
		}
		return
	}

	if appPackage != nil {
		appPackage.SetupRoutes(app)
		if project_name.AppSettings.Debug {
			appLogger.Printf("Routes set up for app: %s", appName)
		}
	} else if project_name.AppSettings.Debug {
		appLogger.Printf("Failed to set up routes for app: %s", appName)
	}
}

func setupRoutes(app *fiber.App) {
	if project_name.AppSettings.Debug {
		appLogger.Println("INSTALLED_APPS:")
		var appNames []string
		for appName := range INSTALLED_APPS {
			appNames = append(appNames, appName)
		}
		sort.Strings(appNames)

		for _, appName := range appNames {
			status := "Disabled"
			if INSTALLED_APPS[appName] {
				status = "Enabled"
			}
			appLogger.Printf("  - %s: %s", appName, status)
		}
	}

	// Monitoring endpoints
	monitor.SetupRoutes(app)

	// Admin routes
	admin.SetupRoutes(app)

	// Set up routes for installed apps
	for appName, isEnabled := range INSTALLED_APPS {
		if isEnabled {
			setupAppRoutes(app, appName)
		}
	}

	// Log all available routes
	logAvailableRoutes(app)
}

func logAvailableRoutes(app *fiber.App) {
	routes := app.Stack()
	for _, route := range routes {
		for _, r := range route {
			appLogger.Printf("Method: %s, Path: %s", r.Method, r.Path)
		}
	}
}
