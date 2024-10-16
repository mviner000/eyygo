package core

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"

	conf "github.com/mviner000/eyymi/src"
	"github.com/mviner000/eyymi/src/config"
	"github.com/mviner000/eyymi/src/constants"
	"github.com/mviner000/eyymi/src/core/decorators"
	"github.com/mviner000/eyymi/src/reverb"
	"github.com/mviner000/eyymi/src/shared"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

var (
	appLogger *log.Logger
	db        *gorm.DB
)

// Define an interface that all apps should implement
type App interface {
	SetupRoutes(app *fiber.App)
}

var jwtSecret = []byte(conf.GetSettings().CSRF.Secret)

type CSRFClaims struct {
	jwt.RegisteredClaims
}

func generateCSRFToken() (string, error) {
	claims := CSRFClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateCSRFToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &CSRFClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func init() {
	if appLogger == nil {
		appLogger = log.New(os.Stdout, "CORE: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Set the default time zone
	loc, err := time.LoadLocation(conf.GetSettings().TimeZone)
	if err != nil {
		appLogger.Fatalf("Invalid time zone: %v", err)
	}
	time.Local = loc

	// Log the time zone if DEBUG is true
	if conf.GetSettings().Debug {
		config.DebugLogf("Time zone set to: %s", conf.GetSettings().TimeZone)
	}

	db, err = gorm.Open(sqlite.Open(conf.GetSettings().Database.Name), &gorm.Config{})
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
}

func ReloadSettings() {
	conf.LoadSettings() // Reload the settings
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

// NewApp initializes and returns a new Fiber application
func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Views:       html.New("./", ".html"),
		ReadTimeout: 5 * time.Second,
	})

	setupMiddleware(app)
	SetupRoutes(app)
	setupCSRFRoutes(app)

	// Set up WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, msg)
			if err != nil {
				break
			}
		}
	}))

	// Log that the WebSocket route is set up
	log.Println(constants.ColorYellow + "WebSocket route set up at /ws" + constants.ColorReset)

	return app
}

func customCORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		postmanToken := c.Get("Postman-Token")

		sharedConfig := shared.GetConfig()
		allowedOrigins := sharedConfig.AllowedOrigins
		isDevelopment := sharedConfig.Environment == "development"

		isAllowedOrigin := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowedOrigin = true
				break
			}
		}

		isPostmanRequest := postmanToken != ""

		if isAllowedOrigin || (isDevelopment && isPostmanRequest) {
			if isDevelopment && isPostmanRequest {
				c.Set("Access-Control-Allow-Origin", "*")
			} else {
				c.Set("Access-Control-Allow-Origin", origin)
			}

			c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Postman-Token")
			c.Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		if !isAllowedOrigin && !(isDevelopment && isPostmanRequest) {
			return fiber.ErrForbidden
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
		TimeZone:   conf.GetSettings().TimeZone,
	}))

	// Use custom CORS middleware
	app.Use(customCORS())

	// CSRF middleware with custom error handler
	// Replace the existing CSRF middleware with our custom one
	app.Use(csrfMiddleware())

	// Custom middlewares
	// app.Use(decorators.RequireHTTPS())
	app.Use(decorators.Logger())
	app.Use(decorators.Throttle(100, 60)) // 100 requests per minute
	app.Use(decorators.DatabaseTransaction(db))
}

func setupCSRFRoutes(app *fiber.App) {
	app.Get("/get-csrf-token", func(c *fiber.Ctx) error {
		token, err := generateCSRFToken()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate CSRF token",
			})
		}
		return c.JSON(fiber.Map{
			"csrf_token": token,
		})
	})
}

func csrfMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "OPTIONS" {
			return c.Next()
		}

		token := c.Get("X-CSRF-Token")
		if token == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "CSRF token not found in request header",
			})
		}

		err := validateCSRFToken(token)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid or expired CSRF token",
			})
		}

		return c.Next()
	}
}

func setupDevelopmentServer() {
	httpPort := os.Getenv("HTTP_PORT")
	wsPort := os.Getenv("WS_PORT")

	if httpPort == "" {
		httpPort = "8000"
	}

	if wsPort == "" {
		wsPort = "3333"
	}

	// Set up HTTP server
	go func() {
		app := fiber.New(fiber.Config{
			Views:       html.New("./", ".html"),
			ReadTimeout: 5 * time.Second,
		})

		setupMiddleware(app)
		reverb.SetupWebSocket(app)
		SetupRoutes(app)

		if conf.GetSettings().Debug {
			appLogger.Printf("Development server started on http://127.0.0.1:%s", httpPort)
		}

		err := app.Listen(":" + httpPort)
		if err != nil {
			appLogger.Fatalf("Failed to start development server: %v", err)
		}
	}()

	// Set up WebSocket server
	go func() {
		app := fiber.New(fiber.Config{
			Views:       html.New("./", ".html"),
			ReadTimeout: 5 * time.Second,
		})

		setupMiddleware(app)
		reverb.SetupWebSocket(app)
		SetupRoutes(app)

		if conf.GetSettings().Debug {
			appLogger.Printf("WebSocket server started on ws://127.0.0.1:%s", wsPort)
		}

		err := app.Listen(":" + wsPort)
		if err != nil {
			appLogger.Fatalf("Failed to start WebSocket server: %v", err)
		}
	}()

	// Block forever
	select {}
}

func setupProductionServer() {
	httpPort := os.Getenv("HTTP_PORT")
	wsPort := os.Getenv("WS_PORT")

	if httpPort == "" {
		httpPort = "8000"
	}

	if wsPort == "" {
		wsPort = "3333"
	}

	// Set up HTTP server
	go func() {
		app := fiber.New(fiber.Config{
			Views:        html.New("./", ".html"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		})

		setupMiddleware(app)
		reverb.SetupWebSocket(app)
		SetupRoutes(app)

		certFile := conf.GetSettings().CertFile
		keyFile := conf.GetSettings().KeyFile

		// Check for wildcard in AllowedOrigins
		for _, origin := range conf.GetSettings().AllowedOrigins {
			if origin == "*" {
				appLogger.Println(constants.ColorRed + "WARNING: Using wildcard '*' in AllowedOrigins in production is not recommended!" + constants.ColorReset)
				break
			}
		}

		if conf.GetSettings().Debug {
			appLogger.Printf("Allowed origins: %v", conf.GetSettings().AllowedOrigins)
		}

		if certFile != "" && keyFile != "" {
			if conf.GetSettings().Debug {
				appLogger.Printf("Starting HTTPS server on port %s", httpPort)
			}
			err := app.ListenTLS(":"+httpPort, certFile, keyFile)
			if err != nil {
				appLogger.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			if conf.GetSettings().Debug {
				appLogger.Printf("Starting HTTP server on port %s", httpPort)
			}
			err := app.Listen(":" + httpPort)
			if err != nil {
				appLogger.Fatalf("Failed to start HTTP server: %v", err)
			}
		}
	}()

	// Set up WebSocket server
	go func() {
		app := fiber.New(fiber.Config{
			Views:        html.New("./", ".html"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		})

		setupMiddleware(app)
		reverb.SetupWebSocket(app)
		SetupRoutes(app)

		certFile := conf.GetSettings().CertFile
		keyFile := conf.GetSettings().KeyFile

		// Check for wildcard in AllowedOrigins
		for _, origin := range conf.GetSettings().AllowedOrigins {
			if origin == "*" {
				appLogger.Println(constants.ColorRed + "WARNING: Using wildcard '*' in AllowedOrigins in production is not recommended!" + constants.ColorReset)
				break
			}
		}

		if conf.GetSettings().Debug {
			appLogger.Printf("Allowed origins: %v", conf.GetSettings().AllowedOrigins)
		}

		if certFile != "" && keyFile != "" {
			if conf.GetSettings().Debug {
				appLogger.Printf("Starting HTTPS WebSocket server on port %s", wsPort)
			}
			err := app.ListenTLS(":"+wsPort, certFile, keyFile)
			if err != nil {
				appLogger.Fatalf("Failed to start HTTPS WebSocket server: %v", err)
			}
		} else {
			if conf.GetSettings().Debug {
				appLogger.Printf("Starting HTTP WebSocket server on port %s", wsPort)
			}
			err := app.Listen(":" + wsPort)
			if err != nil {
				appLogger.Fatalf("Failed to start HTTP WebSocket server: %v", err)
			}
		}
	}()

	// Block forever
	select {}
}
