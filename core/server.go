package core

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/mviner000/eyymi/admin"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/core/decorators"
	"github.com/mviner000/eyymi/monitor"
	"github.com/mviner000/eyymi/reverb"
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
	reverb.SetLogger(appLogger)

	var err error
	db, err = gorm.Open(sqlite.Open(config.GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
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
	appPackage, err := getAppPackage(appName)
	if err != nil {
		if config.AppSettings.Debug {
			appLogger.Printf("Error setting up app: %v", err)
		}
		return
	}

	if appPackage != nil {
		appPackage.SetupRoutes(app)
		if config.AppSettings.Debug {
			appLogger.Printf("Routes set up for app: %s", appName)
		}
	} else if config.AppSettings.Debug {
		appLogger.Printf("Failed to set up routes for app: %s", appName)
	}
}

func setupRoutes(app *fiber.App) {
	if config.AppSettings.Debug {
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

func setupMiddleware(app *fiber.App, isProd bool) {
	// Recover middleware
	app.Use(recover.New())

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// CORS middleware
	corsConfig := cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}
	if isProd {
		corsConfig.AllowOrigins = config.GetAllowedOrigins()
	} else {
		corsConfig.AllowOrigins = "*"
	}
	app.Use(cors.New(corsConfig))

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

	setupMiddleware(app, false)
	reverb.SetupWebSocket(app)
	setupRoutes(app)

	wsPort := config.GetWebSocketPort()

	if config.AppSettings.Debug {
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

	setupMiddleware(app, true)
	reverb.SetupWebSocket(app)
	setupRoutes(app)

	wsPort := config.GetWebSocketPort()
	certFile := config.GetCertFile()
	keyFile := config.GetKeyFile()

	if config.AppSettings.Debug {
		appLogger.Printf("Allowed origins: %s", config.GetAllowedOrigins())
	}

	if certFile != "" && keyFile != "" {
		if config.AppSettings.Debug {
			appLogger.Printf("Starting HTTPS server on port %s", wsPort)
		}
		err := app.ListenTLS(":"+wsPort, certFile, keyFile)
		if err != nil {
			appLogger.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		if config.AppSettings.Debug {
			appLogger.Printf("Starting HTTP server on port %s", wsPort)
		}
		err := app.Listen(":" + wsPort)
		if err != nil {
			appLogger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}
