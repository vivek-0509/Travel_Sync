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

// Redirect to Google login
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	state := generateRandomState()
	// Short-lived state cookie for OAuth verification
	c.SetCookie("oauth_state", state, 300, "/", "", true, true) // Secure=true for production

	// Add prompt=consent to force consent screen every time
	url := h.CustomOAuth2Service.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback after Google login
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	cookieState, _ := c.Cookie("oauth_state")
	if state == "" || cookieState == "" || state != cookieState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	jwtToken, created, user, err := h.CustomOAuth2Service.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set JWT cookie for cross-origin requests
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://www.travelsync.space" // updated fallback
	}

	// For cross-origin cookies, don't set domain - let browser handle it
	// Setting domain can cause issues with cross-origin requests
	c.Header("Set-Cookie", fmt.Sprintf(
		"jwt_token=%s; Path=/; Max-Age=%d; HttpOnly; Secure; SameSite=None",
		jwtToken,
		3600*24*8,
	))

	// Check if profile is complete and redirect accordingly
	redirectURL := frontendURL + "/auth/success"
	if created || !h.CustomOAuth2Service.AuthService.IsProfileComplete(user) {
		redirectURL = frontendURL + "/auth/success/?new=1"
	}
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// Logout handler
func (h *OAuthHandler) Logout(c *gin.Context) {
	// Try to get JWT claims to extract access and refresh tokens for revocation
	if jwtClaims, exists := c.Get("jwt_claims"); exists {
		if claims, ok := jwtClaims.(*service.CustomClaims); ok {
			// Revoke both Google OAuth access and refresh tokens
			if err := h.CustomOAuth2Service.RevokeGoogleToken(c.Request.Context(), claims.AccessToken, claims.RefreshToken); err != nil {
				// Log error but don't fail logout - token might already be expired/revoked
				fmt.Printf("Warning: Failed to revoke Google OAuth tokens: %v\n", err)
			}
		}
	}

	// Clear JWT cookie with same settings as when it was set
	c.Header("Set-Cookie", "jwt_token=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=None")
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

// generateRandomState generates a random string for OAuth state
func generateRandomState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
