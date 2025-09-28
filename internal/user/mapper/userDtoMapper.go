package mapper

import (
	"Travel_Sync/internal/user/entity"
	"Travel_Sync/internal/user/models"
)

func FromUserCreateDto(createDto *models.UserCreateDto) *entity.User {
	return &entity.User{
		Email: createDto.Email,
		Batch: createDto.Batch,
	}
}

func FromUserUpdateDto(updateDto *models.UserUpdateDto, user *entity.User) *entity.User {
	if updateDto.Name != "" {
		user.Name = updateDto.Name
	}
	if updateDto.PhoneNumber != "" {
		user.PhoneNumber = updateDto.PhoneNumber
	}
	return user
}
