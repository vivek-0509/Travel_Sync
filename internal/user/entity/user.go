package entity

import "time"

type User struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;not null"`
	Name        string    `gorm:"size:255;not null"`
	Email       string    `gorm:"not null;unique"`
	Batch       int       `gorm:"not null" `
	PhoneNumber int       `gorm:"not null;size:10" `
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
