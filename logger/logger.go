package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

type Logger struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		InfoLogger:    log.New(os.Stdout, "", 0),
		WarningLogger: log.New(os.Stdout, "", 0),
		ErrorLogger:   log.New(os.Stderr, "", 0),
	}
}

// PrintServerStatus prints the server startup information similar to Django
func (l *Logger) PrintServerStatus(host string, port string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	l.InfoLogger.Printf("\nStarting development server at %s", time.Now().Format("2006-01-02 15:04:05"))
	l.InfoLogger.Printf("Watching for file changes with server reload")
	l.InfoLogger.Printf("%s %s", yellow("Eytgo version"), "Fiber-Go/1.0")
	l.InfoLogger.Printf("%s %s", yellow("Operating System"), os.Getenv("OS"))
	l.InfoLogger.Printf("\nRunning server on: %s", green(fmt.Sprintf("http://%s:%s/", host, port)))
	l.InfoLogger.Printf("%s %s", cyan("Quit the server with"), "CONTROL-C\n")
}

// RequestLogger middleware for logging HTTP requests
func RequestLogger() fiber.Handler {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate processing time
		duration := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Color-code the status
		var statusColored string
		switch {
		case status >= 500:
			statusColored = red(status)
		case status >= 400:
			statusColored = yellow(status)
		case status >= 300:
			statusColored = cyan(status)
		default:
			statusColored = green(status)
		}

		// Format the log message similar to Django
		message := fmt.Sprintf("[%s] %s %s %s %s %s %s",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Method(),
			statusColored,
			c.Path(),
			duration.Round(time.Millisecond),
			cyan(c.IP()),
			yellow(c.Get("User-Agent")),
		)

		log.Println(message)

		return err
	}
}

// ErrorLogger middleware for logging errors
func ErrorLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			red := color.New(color.FgRed).SprintFunc()
			log.Printf("%s Error: %v", red("[ERROR]"), err)
		}

		return err
	}
}
