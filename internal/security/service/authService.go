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
// Returns: user, created(bool), error
func (authService *AuthService) GetOrCreateUser(email string) (*entity.User, bool, error) {
    user, err := authService.UserService.GetUserByEmail(email)

    if err != nil || user == nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            user, err = authService.UserService.CreateUser(email)
            if err != nil {
                return nil, false, err
            }
            return user, true, nil
        }
        return nil, false, err
    }
    return user, false, nil
}

// IsProfileComplete checks if user has completed their profile (name and phone number)
func (authService *AuthService) IsProfileComplete(user *entity.User) bool {
    return user.Name != "" && user.PhoneNumber != ""
}