package reverb

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	conf "github.com/mviner000/eyygo/src"
)

var (
	logger *log.Logger
)

func SetLogger(l *log.Logger) {
	if l == nil {
		logger = log.New(os.Stdout, "WEBSOCKET: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = l
	}
}

func SetupWebSocket(app *fiber.App) {
	if logger == nil {
		SetLogger(nil) // This will create a default logger
	}

	wsPort := os.Getenv("WS_PORT")
	if wsPort == "" {
		wsPort = "3000" // Default to 3000 if not set
	}

	// Convert wsPort to int
	wsPortInt, err := strconv.Atoi(wsPort)
	if err != nil {
		logger.Fatalf("Invalid WS_PORT: %v", err)
	}

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(handleWebSocket))

	// Log WebSocket configuration if DEBUG is true
	if conf.GetSettings().Debug {
		logger.Printf("WebSocket server configured on port %d", wsPortInt)
	}
}

func handleWebSocket(c *websocket.Conn) {
	if logger == nil {
		SetLogger(nil) // This will create a default logger if it's still nil
	}

	// Basic WebSocket connection handler
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Printf("Error reading message: %v", err)
			}
			return
		}

		if messageType == websocket.TextMessage {
			logger.Printf("Received message: %s", string(message))
			// Handle the message as needed
		}
	}
}
