package service

import (
	"Travel_Sync/internal/user/service"
	"errors"

	"gorm.io/gorm"
)

type AuthService struct {
	UserService *service.UserService
}

func NewAuthService(userService *service.UserService) *AuthService {
	return &AuthService{UserService: userService}
}

// Get user by email or create a new user of not exists
func (authService *AuthService) GetOrCreateUser(email string) error {
	user, err := authService.UserService.GetUserByEmail(email)

	if err != nil || user == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			authService.UserService.CreateUser(email)
			return nil
		}
		return err
	}
	return nil
}
