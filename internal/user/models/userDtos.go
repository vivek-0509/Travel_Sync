package models

type UserCreateDto struct {
	Email string `validate:"required,email"`
	Batch int    `validate:"required,gte=1"`
}

type UserUpdateDto struct {
	Name        string `validate:"max=255"`
	PhoneNumber *int   `validate:"phone"`
}
