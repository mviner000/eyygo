package core

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/mviner000/eyymi/app_name"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/admin"
	"github.com/mviner000/eyymi/eyygo/constants"
	"github.com/mviner000/eyymi/eyygo/core/decorators"
	"github.com/mviner000/eyymi/eyygo/monitor"
	"github.com/mviner000/eyymi/eyygo/reverb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	appLogger *log.Logger
	db        *gorm.DB
)

// Define an interface that all apps should implement
type App interface {
	SetupRoutes(app *fiber.App)
}

func init() {
	if appLogger == nil {
		appLogger = log.New(os.Stdout, "CORE: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Set the default time zone
	loc, err := time.LoadLocation(app_name.AppSettings.TimeZone)
	if err != nil {
		appLogger.Fatalf("Invalid time zone: %v", err)
	}
	time.Local = loc

	// Add this debug logging
	config.DebugLogf("Time zone set to: %s", app_name.AppSettings.TimeZone)

	db, err = gorm.Open(sqlite.Open(app_name.AppSettings.Database.Name), &gorm.Config{})
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
}

func ReloadSettings() {
	app_name.LoadSettings() // Reload settings using the function from app_name
	log.Println("Settings reloaded")
}

func RunCommand() {
	ReloadSettings() // Ensure settings are reloaded at the start

	nodeEnv := os.Getenv("NODE_ENV")
	isProduction := nodeEnv == "production"

	if isProduction {
		config.DebugLogf("Running in production mode")
		setupProductionServer()
	} else {
		config.DebugLogf("Running in development mode")
		setupDevelopmentServer()
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
	appPackage, err := getAppPackage(appName)
	if err != nil {
		if app_name.AppSettings.Debug {
			appLogger.Printf("Error setting up app: %v", err)
		}
		return
	}

	if appPackage != nil {
		appPackage.SetupRoutes(app)
		if app_name.AppSettings.Debug {
			appLogger.Printf("Routes set up for app: %s", appName)
		}
	} else if app_name.AppSettings.Debug {
		appLogger.Printf("Failed to set up routes for app: %s", appName)
	}
}

func setupRoutes(app *fiber.App) {
	if app_name.AppSettings.Debug {
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
}

func customCORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

func setupMiddleware(app *fiber.App) {
	// Recover middleware
	app.Use(recover.New())

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   app_name.AppSettings.TimeZone,
	}))

	// Use custom CORS middleware
	app.Use(customCORS())

	// Custom middlewares
	// app.Use(decorators.RequireHTTPS())
	app.Use(decorators.Logger())
	app.Use(decorators.Throttle(100, 60)) // 100 requests per minute
	app.Use(decorators.DatabaseTransaction(db))
}

func setupDevelopmentServer() {
	app := fiber.New(fiber.Config{
		Views:       html.New("./", ".html"),
		ReadTimeout: 5 * time.Second,
	})

	setupMiddleware(app) // Removed the second argument
	reverb.SetupWebSocket(app)
	setupRoutes(app)

	wsPort := app_name.AppSettings.WebSocket.Port

	if app_name.AppSettings.Debug {
		appLogger.Printf("Development server started on http://127.0.0.1:%s", wsPort)
	}

	err := app.Listen(":" + wsPort)
	if err != nil {
		appLogger.Fatalf("Failed to start development server: %v", err)
	}
}

func setupProductionServer() {
	app := fiber.New(fiber.Config{
		Views:        html.New("./", ".html"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	setupMiddleware(app)
	reverb.SetupWebSocket(app)
	setupRoutes(app)

	wsPort := app_name.AppSettings.WebSocket.Port
	certFile := app_name.AppSettings.CertFile
	keyFile := app_name.AppSettings.KeyFile

	// Check for wildcard in AllowedOrigins
	for _, origin := range app_name.AppSettings.AllowedOrigins {
		if origin == "*" {
			appLogger.Println(constants.ColorRed + "WARNING: Using wildcard '*' in AllowedOrigins in production is not recommended!" + constants.ColorReset)
			break
		}
	}

	if app_name.AppSettings.Debug {
		appLogger.Printf("Allowed origins: %v", app_name.AppSettings.AllowedOrigins)
	}

	if certFile != "" && keyFile != "" {
		if app_name.AppSettings.Debug {
			appLogger.Printf("Starting HTTPS server on port %s", wsPort)
		}
		err := app.ListenTLS(":"+wsPort, certFile, keyFile)
		if err != nil {
			appLogger.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		if app_name.AppSettings.Debug {
			appLogger.Printf("Starting HTTP server on port %s", wsPort)
		}
		err := app.Listen(":" + wsPort)
		if err != nil {
			appLogger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}
