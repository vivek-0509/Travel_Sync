package models

type TravelTicketCreateDto struct {
	Source       string `json:"source" binding:"required"`
	Destination  string `json:"destination" binding:"required"`
	DepartureAt  string `json:"departure_at" binding:"required"` // RFC3339 e.g., 2025-10-01T14:30:00Z
	TimeDiffMins int    `json:"time_diff_mins" binding:"required,min=0,max=720"`
	EmptySeats   int    `json:"empty_seats" binding:"required,min=1,max=10"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
}

type TravelTicketUpdateDto struct {
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	DepartureAt  string `json:"departure_at"` // RFC3339, optional
	TimeDiffMins int    `json:"time_diff_mins"`
	EmptySeats   int    `json:"empty_seats"`
	PhoneNumber  string `json:"phone_number"`
	Status       string `json:"status"` // "open" or "closed"
}

type TravelTicketUserResponseDto struct {
	ID           int64  `json:"id"`
	StudentName  string `json:"student_name"`
	StudentBatch string `json:"student_batch"`
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	Date         string `json:"date"` // 2006-01-02
	Time         string `json:"time"` // 15:04
	EmptySeats   int    `json:"empty_seats"`
	PhoneNumber  string `json:"phone_number"`
}
