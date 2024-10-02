package decorators

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// GzipMiddleware applies gzip compression to responses
func GzipMiddleware() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelDefault,
	})
}
