package service

import (
	"Travel_Sync/internal/user/entity"
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

// Get user by email or create a new user if not exists
func (authService *AuthService) GetOrCreateUser(email string) (*entity.User, error) {
	user, err := authService.UserService.GetUserByEmail(email)

	if err != nil || user == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = authService.UserService.CreateUser(email)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	return user, nil
}
