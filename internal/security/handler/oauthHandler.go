package handler

import (
	"Travel_Sync/internal/security/service"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	CustomOAuth2Service *service.CustomOAuth2Service
}

func NewOAuthHandler(customOAuth2Service *service.CustomOAuth2Service) *OAuthHandler {
	return &OAuthHandler{CustomOAuth2Service: customOAuth2Service}
}

// Redirect to google
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	state := generateRandomState() // implement secure random
	// set short-lived cookie with state
	c.SetCookie("oauth_state", state, 300, "/", "", false, true) // secure=true in prod
	url := h.CustomOAuth2Service.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {

	// Verify state
	state := c.Query("state")
	cookieState, _ := c.Cookie("oauth_state")
	if state == "" || cookieState == "" || state != cookieState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	code := c.Query("code")

	jwtToken, created, err := h.CustomOAuth2Service.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//Set Jwt as HTTP-only cookie
	// Configure cookie for production readiness
	secure := os.Getenv("COOKIE_SECURE") == "true"
	domain := os.Getenv("COOKIE_DOMAIN")
	c.SetCookie(
		"jwt_token",
		jwtToken,
		3600*24*8,
		"/",
		domain,
		secure,
		true,
	)

	// redirect to frontend app success page (prefer env var)
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // dev fallback
	}
	redirectURL := frontendURL + "/auth/success"
	if created {
		redirectURL += "?new=1"
	}
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)

}

// Logout handler
func (h *OAuthHandler) Logout(c *gin.Context) {
	// Clear the JWT cookie
	c.SetCookie(
		"jwt_token", // cookie name
		"",          // empty value
		-1,          // max-age (negative means delete immediately)
		"/",         // path
		"",          // domain
		false,       // secure
		true,        // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Get current user info
func (h *OAuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEmail, _ := c.Get("user_email")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   userEmail,
	})
}

func generateRandomState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
