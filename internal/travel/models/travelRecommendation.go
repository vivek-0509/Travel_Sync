package models

import tentity "Travel_Sync/internal/travel/entity"

type ScoredTicket struct {
	Ticket tentity.TravelTicket `json:"ticket"`
	Score  float64              `json:"score"`
	Date   string               `json:"date"`
	Time   string               `json:"time"`
}

type RecommendationResult struct {
	BestMatch         *ScoredTicket  `json:"best_match"`
	BestGroup         []ScoredTicket `json:"best_group"`
	OtherAlternatives []ScoredTicket `json:"other_alternatives"`
}
