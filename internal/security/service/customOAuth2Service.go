package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	//Generate Jwt Token with user ID
    jwtToken, err := service.JWTService.CreateJWT(user.ID, googleUser.Email)
	if err != nil {
        return "", false, err
	}

    return jwtToken, created, nil
}

func ExtractDomain(email string) string {
	parts := strings.Split(email, "@")
	return parts[1]
}
