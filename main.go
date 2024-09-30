package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/mviner000/eyymi/monitor"
	"github.com/mviner000/eyymi/reverb"
)

var (
	logger *log.Logger
)

func init() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	reverb.SetLogger(logger)
}

func main() {
	if os.Getenv("NODE_ENV") == "development" {
		setupDevelopmentServer()
	} else {
		setupProductionServer()
	}
}

// Common route setup for both dev and production
func setupRoutes(app *fiber.App) {
	// Monitoring endpoints
	app.Get("/status", monitor.HandleStatus)
	app.Get("/status/server-info", monitor.HandleStatus)
	app.Get("/status/cpu", monitor.HandleStatus)
	app.Get("/status/ram", monitor.HandleStatus)
	app.Get("/status/storage", monitor.HandleStatus)
	app.Get("/status/old", monitor.HandleStatus)
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

	// Get WebSocket port from environment variable
	wsPort := os.Getenv("WS_PORT")
	if wsPort == "" {
		wsPort = "3000" // Default to 3000 if not set
	}

	// Log WebSocket server start
	logger.Printf("WebSocket server started on http://127.0.0.1:%s", wsPort)

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

	// Get allowed origins from environment variable
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "https://eyymi.site"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	reverb.SetupWebSocket(app)
	setupRoutes(app) // Call common route setup

	port := os.Getenv("WS_PORT")
	if port == "" {
		port = "3000"
	}

	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")

	logger.Printf("Allowed origins: %s", allowedOrigins)

	if certFile != "" && keyFile != "" {
		logger.Printf("Starting HTTPS server on port %s", port)
		err := app.ListenTLS(":"+port, certFile, keyFile)
		if err != nil {
			logger.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		logger.Printf("Starting HTTP server on port %s", port)
		err := app.Listen(":" + port)
		if err != nil {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}
