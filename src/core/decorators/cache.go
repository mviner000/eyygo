package decorators

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

var (
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c = cache.New(5*time.Minute, 10*time.Minute)
)

// CachePageDecorator caches the response of a Fiber handler
func CachePageDecorator(duration time.Duration) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		key := generateCacheKey(ctx)

		// Try to get the cached response
		if cachedResponse, found := c.Get(key); found {
			return ctx.Send(cachedResponse.([]byte))
		}

		// If not found in cache, call the next handler
		err := ctx.Next()
		if err != nil {
			return err
		}

		// Cache the response
		c.Set(key, ctx.Response().Body(), duration)

		return nil
	}
}

// generateCacheKey creates a unique key for each request
func generateCacheKey(ctx *fiber.Ctx) string {
	// You might want to customize this based on your needs
	data := ctx.Method() + ctx.OriginalURL() + string(ctx.Request().Body())
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ClearCache clears the entire cache
func ClearCache() {
	c.Flush()
}

// ClearCacheKey removes a specific key from the cache
func ClearCacheKey(key string) {
	c.Delete(key)
}
