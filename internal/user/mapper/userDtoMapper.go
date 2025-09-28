package mapper

import (
	"Travel_Sync/internal/user/entity"
	"Travel_Sync/internal/user/models"
	"strings"
)

func FromUserEmail(email string) *entity.User {
	return &entity.User{
		Email: email,
		Batch: ExtractBatch(email),
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

func ExtractBatch(email string) string {
	parts := strings.Split(email, ".")
	userNamePart := parts[1]

	yearSuffix := userNamePart[:2]
	batchStr := "Batch20" + yearSuffix
	return batchStr
}
