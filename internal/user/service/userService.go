package service

import (
	"Travel_Sync/internal/user/entity"
	"Travel_Sync/internal/user/mapper"
	"Travel_Sync/internal/user/models"
	"Travel_Sync/internal/user/repository"
)

type UserService struct {
	Repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (svc *UserService) CreateUser(createDto *models.UserCreateDto) (*entity.User, error) {
	user := mapper.FromUserCreateDto(createDto)
	user, err := svc.Repo.Create(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *UserService) GetByID(userID int) (*entity.User, error) {
	user, err := svc.Repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *UserService) GetAll() ([]entity.User, error) {
	users, err := svc.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (svc *UserService) DeleteByID(userID int) error {
	err := svc.Repo.Delete(userID)
	if err != nil {
		return err
	}
	return nil
}

func (svc *UserService) UpdateUser(userId int, updateDto *models.UserUpdateDto) error {
	user, err := svc.Repo.GetByID(userId)
	if err != nil {
		return err
	}
	user = mapper.FromUserUpdateDto(updateDto, user)

	err = svc.Repo.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}
