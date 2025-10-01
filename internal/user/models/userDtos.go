package models

type UserCreateDto struct {
	Email string `validate:"required,email"`
	Batch string `validate:"required,gte=1"`
}

type UserUpdateDto struct {
    Name        string `json:"name" validate:"max=255"`
    PhoneNumber string `json:"phone_number" validate:"phone"`
}
