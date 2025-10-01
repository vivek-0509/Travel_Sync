package models

import (
	"time"
)

type MinimalUser struct {
	Name  string `json:"name"`
	Batch string `json:"batch"`
	Email string `json:"email"`
}

// PublicTicket is a redacted view of TravelTicket for recommendations
// that hides sensitive identifiers like id and user_id.
type PublicTicket struct {
	Source       string    `json:"source"`
	Destination  string    `json:"destination"`
	EmptySeats   int       `json:"empty_seats"`
	DepartureAt  time.Time `json:"departure_at"`
	TimeDiffMins int       `json:"time_diff_mins"`
	PhoneNumber  string    `json:"phone_number"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ScoredTicket struct {
	Ticket      PublicTicket `json:"ticket"`
	Score       float64      `json:"score"`
	Date        string       `json:"date"`
	Time        string       `json:"time"`
	User        MinimalUser  `json:"user"`
	CandidateID int64        `json:"-"` // internal use only, not exposed
}

type RecommendationResult struct {
	BestMatch         *ScoredTicket  `json:"best_match"`
	BestGroup         []ScoredTicket `json:"best_group"`
	OtherAlternatives []ScoredTicket `json:"other_alternatives"`
}
