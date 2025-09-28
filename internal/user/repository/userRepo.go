package repository

import (
	"Travel_Sync/internal/user/entity"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// Create a new user
func (r *UserRepo) Create(user *entity.User) (*entity.User, error) {
	err := r.DB.Create(&user).Error
	return user, err
}

// GetById Get User by ID
func (r *UserRepo) GetByID(userID int64) (*entity.User, error) {
	var user entity.User
	err := r.DB.First(&user, userID).Error
	return &user, err
}

// GetAll User
func (r *UserRepo) GetAll() ([]entity.User, error) {
	var users []entity.User
	err := r.DB.Find(&users).Error
	return users, err
}

// Update User
func (r *UserRepo) UpdateUser(user *entity.User) (*entity.User, error) {
	err := r.DB.Save(user).Error
	return user, err
}

// Delete User By ID
func (r *UserRepo) Delete(userID int64) error {
	return r.DB.Delete(&entity.User{ID: userID}).Error
}
