package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// RateLimitConfig defines configuration for each limiter
type RateLimitConfig struct {
	Rate   rate.Limit // requests per second
	Burst  int        // burst capacity
	Prefix string     // key prefix
}

// Global limiter map
var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// Rate limiters for specific endpoints
var (
	AuthRateLimit = RateLimitConfig{
		Rate:   rate.Every(time.Minute / 30), // 30 requests per minute
		Burst:  60,
		Prefix: "auth",
	}

	GeneralRateLimit = RateLimitConfig{
		Rate:   rate.Every(time.Hour / 1200), // 1200 requests per hour
		Burst:  1800,
		Prefix: "general",
	}

	RecommendationRateLimit = RateLimitConfig{
		Rate:   rate.Every(time.Minute / 60), // 60 requests per minute
		Burst:  120,
		Prefix: "recommend",
	}
)

// getLimiter returns (or creates) a limiter for the given key
func getLimiter(key string, config RateLimitConfig) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[key]
	if !exists {
		limiter = rate.NewLimiter(config.Rate, config.Burst)
		limiters[key] = limiter
	}
	return limiter
}

// RateLimitMiddleware creates a Gin middleware for a given config
func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		// 1. If logged in → use user_id
		if uid, exists := c.Get("user_id"); exists {
			key = config.Prefix + ":user:" + uid.(string)
		} else {
			// 2. Guests → assign unique guest_id cookie
			guestID, err := c.Cookie("guest_id")
			if err != nil || guestID == "" {
				guestID = uuid.New().String()
				// Set cookie for 1h
				c.SetCookie("guest_id", guestID, 3600, "/", "", false, true)
			}
			key = config.Prefix + ":guest:" + guestID

			// 3. Fallback → IP + UA
			if guestID == "" {
				ua := c.Request.UserAgent()
				key = config.Prefix + ":ipua:" + c.ClientIP() + ":" + ua
			}
		}

		limiter := getLimiter(key, config)

		// Deny if not allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":     false,
				"error":       "Rate limit exceeded. Please try again later.",
				"retry_after": limiter.Reserve().Delay().Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Specific middlewares
func AuthRateLimiter() gin.HandlerFunc {
	return RateLimitMiddleware(AuthRateLimit)
}

func GeneralRateLimiter() gin.HandlerFunc {
	return RateLimitMiddleware(GeneralRateLimit)
}

func RecommendationRateLimiter() gin.HandlerFunc {
	return RateLimitMiddleware(RecommendationRateLimit)
}
