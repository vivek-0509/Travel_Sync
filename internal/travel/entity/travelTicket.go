package entity

import "time"

type TravelTicket struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Source       string    `gorm:"size:255;not null" json:"source"`
	Destination  string    `gorm:"size:255;not null" json:"destination"`
	EmptySeats   int       `gorm:"not null" json:"empty_seats"`
	DepartureAt  time.Time `gorm:"type:timestamptz;not null" json:"departure_at"`
	TimeDiffMins int       `gorm:"not null" json:"time_diff_mins"`
	UserID       int64     `gorm:"not null" json:"user_id"`
	PhoneNumber  string    `gorm:"size:15;not null" json:"phone_number"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
