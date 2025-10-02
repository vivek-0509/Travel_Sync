package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type RateLimitConfig struct {
	Rate   string // e.g., "10-M" for 10 requests per minute
	Burst  int    // burst limit
	Prefix string // key prefix
}

var (
	// Default rate limiters
	AuthRateLimit = RateLimitConfig{
		Rate:   "30-M", // 5 requests per minute for auth endpoints
		Burst:  60,
		Prefix: "auth",
	}

	GeneralRateLimit = RateLimitConfig{
		Rate:   "1200-H", // 100 requests per hour for general endpoints
		Burst:  1800,
		Prefix: "general",
	}

	RecommendationRateLimit = RateLimitConfig{
		Rate:   "60-M", // 20 requests per minute for recommendations
		Burst:  120,
		Prefix: "recommend",
	}
)

func RateLimiter(config RateLimitConfig) gin.HandlerFunc {
	// Create rate limiter
	rate, err := limiter.NewRateFromFormatted(config.Rate)
	if err != nil {
		panic(err)
	}

	// Create memory store
	store := memory.NewStore()

	// Create limiter instance
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get client identifier
		clientIP := c.ClientIP()

		// Create context
		ctx := context.Background()

		// Create key with prefix: prefer per-user when authenticated, else fall back to IP
		key := ""
		if userIDVal, exists := c.Get("user_id"); exists {
			switch id := userIDVal.(type) {
			case int64:
				key = config.Prefix + ":uid:" + strconv.FormatInt(id, 10)
			case int:
				key = config.Prefix + ":uid:" + strconv.FormatInt(int64(id), 10)
			case string:
				key = config.Prefix + ":uid:" + id
			default:
				key = config.Prefix + ":" + clientIP
			}
		} else {
			key = config.Prefix + ":" + clientIP
		}

		// Check rate limit
		context, err := instance.Get(ctx, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Rate limit check failed",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Check if limit exceeded
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":     false,
				"error":       "Rate limit exceeded. Please try again later.",
				"retry_after": context.Reset,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Specific rate limiters for different endpoints
func AuthRateLimiter() gin.HandlerFunc {
	return RateLimiter(AuthRateLimit)
}

func GeneralRateLimiter() gin.HandlerFunc {
	return RateLimiter(GeneralRateLimit)
}

func RecommendationRateLimiter() gin.HandlerFunc {
	return RateLimiter(RecommendationRateLimit)
}
