package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type CustomOAuth2Service struct {
	OAuthConfig *oauth2.Config
	AuthService *AuthService
	JWTService  *JWTService
}

func NewCustomOAuth2Service(oauthConfig *oauth2.Config, authService *AuthService, jwtService *JWTService) *CustomOAuth2Service {
	return &CustomOAuth2Service{
		OAuthConfig: oauthConfig,
		AuthService: authService,
		JWTService:  jwtService,
	}
}

// GoogleCallback exchanges the code, ensures user exists, and returns jwt + created flag
func (service *CustomOAuth2Service) GoogleCallback(ctx context.Context, code string) (string, bool, error) {

	if code == "" {
        return "", false, fmt.Errorf("code not found")
	}

	//Exchanging code for token
    token, err := service.OAuthConfig.Exchange(ctx, code)
	if err != nil {
        return "", false, fmt.Errorf("token exchange failed: %w", err)
	}

	//Fetch User Info from google
	client := service.OAuthConfig.Client(ctx, token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
        return "", false, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
        return "", false, err
	}

	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

    if err := json.Unmarshal(body, &googleUser); err != nil {
        return "", false, fmt.Errorf("failed to parse user info: %w", err)
	}

	if ExtractDomain(googleUser.Email) != "sst.scaler.com" {
        return "", false, errors.New("login with your Scaler Student Email")
	}

    user, created, err := service.AuthService.GetOrCreateUser(googleUser.Email)
	if err != nil {
        return "", false, err
	}

	//Generate Jwt Token with user ID and access token
    jwtToken, err := service.JWTService.CreateJWT(user.ID, googleUser.Email, token.AccessToken)
	if err != nil {
        return "", false, err
	}

    return jwtToken, created, nil
}

// RevokeGoogleToken revokes the Google OAuth access token
func (service *CustomOAuth2Service) RevokeGoogleToken(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("access token is empty")
	}

	// Google's token revocation endpoint
	revokeURL := "https://oauth2.googleapis.com/revoke"
	
	// Create HTTP request to revoke token
	req, err := http.NewRequestWithContext(ctx, "POST", revokeURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create revoke request: %w", err)
	}

	// Add token as form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = map[string][]string{
		"token": {accessToken},
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token revocation failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func ExtractDomain(email string) string {
	parts := strings.Split(email, "@")
	return parts[1]
}
