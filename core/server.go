package core

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/mviner000/eyymi/admin"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/monitor"
	"github.com/mviner000/eyymi/reverb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	logger *log.Logger
	db     *gorm.DB
)

// Define an interface that all apps should implement
type App interface {
	SetupRoutes(app *fiber.App)
}

func init() {
	if logger == nil {
		logger = log.New(os.Stdout, "CORE: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	reverb.SetLogger(logger)

	var err error
	db, err = gorm.Open(sqlite.Open(config.GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
}

func RunCommand() {
	if config.IsDevelopment() {
		setupDevelopmentServer()
	} else {
		setupProductionServer()
	}
}

var exampleapp App = nil

// Mocked INSTALLED_APPS-like structure
var INSTALLED_APPS = map[string]bool{
	"exampleapp": true,  // This app is enabled
	"otherapp":   false, // This app is disabled and won't be set up
	// Add more apps as needed
}

func getAppPackage(appName string) (App, error) {
	switch appName {
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
	// Get the app package statically using the new helper function
	appPackage, err := getAppPackage(appName)
	if err != nil {
		if config.AppSettings.Debug {
			logger.Printf("Error setting up app: %v", err)
		}
		return
	}

	// If appPackage is found, call its SetupRoutes function
	if appPackage != nil {
		appPackage.SetupRoutes(app)
		if config.AppSettings.Debug {
			logger.Printf("Routes set up for app: %s", appName)
		}
	} else if config.AppSettings.Debug {
		logger.Printf("Failed to set up routes for app: %s", appName)
	}
}

func setupRoutes(app *fiber.App) {
	// Debug logging for INSTALLED_APPS
	if config.AppSettings.Debug {
		logger.Println("INSTALLED_APPS:")
		// Create a sorted list of app names for consistent output
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
			logger.Printf("  - %s: %s", appName, status)
		}
	}

	// Monitoring endpoints
	monitor.SetupRoutes(app)

	// Admin routes
	_, adminSetup := admin.NewAdminModule(db)
	adminSetup(app)

	// Set up routes for installed apps
	for appName, isEnabled := range INSTALLED_APPS {
		if isEnabled {
			setupAppRoutes(app, appName)
		}
	}
}

func setupDevelopmentServer() {
	app := fiber.New(fiber.Config{
		Views: html.New("./", ".html"),
	})

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	reverb.SetupWebSocket(app)
	setupRoutes(app) // Call common route setup

	wsPort := config.GetWebSocketPort()

	// Log WebSocket server start only if debug is true
	if config.AppSettings.Debug {
		logger.Printf("WebSocket server started on http://127.0.0.1:%s", wsPort)
	}

	// Start the server
	err := app.Listen(":" + wsPort)
	if err != nil {
		logger.Fatalf("Failed to start WebSocket server: %v", err)
	}
}

func setupProductionServer() {
	app := fiber.New(fiber.Config{
		Views: html.New("./", ".html"),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetAllowedOrigins(),
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	reverb.SetupWebSocket(app)
	setupRoutes(app) // Call common route setup

	wsPort := config.GetWebSocketPort()
	certFile := config.GetCertFile()
	keyFile := config.GetKeyFile()

	if config.AppSettings.Debug {
		logger.Printf("Allowed origins: %s", config.GetAllowedOrigins())
	}

	if certFile != "" && keyFile != "" {
		if config.AppSettings.Debug {
			logger.Printf("Starting HTTPS server on port %s", wsPort)
		}
		err := app.ListenTLS(":"+wsPort, certFile, keyFile)
		if err != nil {
			logger.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		if config.AppSettings.Debug {
			logger.Printf("Starting HTTP server on port %s", wsPort)
		}
		err := app.Listen(":" + wsPort)
		if err != nil {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}
