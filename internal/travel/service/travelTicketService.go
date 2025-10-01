package service

import (
	tentity "Travel_Sync/internal/travel/entity"
	"Travel_Sync/internal/travel/mapper"
	"Travel_Sync/internal/travel/models"
	"Travel_Sync/internal/travel/repository"
	urepo "Travel_Sync/internal/user/repository"
	"errors"
	"math"
	"sort"
	"time"
)

type TravelTicketService struct {
	Repo     *repository.TravelTicketRepo
	UserRepo *urepo.UserRepo
}

func NewTravelTicketService(repo *repository.TravelTicketRepo, userRepo *urepo.UserRepo) *TravelTicketService {
	return &TravelTicketService{Repo: repo, UserRepo: userRepo}
}

func (s *TravelTicketService) Create(userID int64, dto *models.TravelTicketCreateDto) (*tentity.TravelTicket, error) {
	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.PhoneNumber == "" && dto.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	ticket, err := mapper.FromCreateDtoToEntity(dto, user)
	if err != nil {
		return nil, err
	}
	if ticket.PhoneNumber == "" {
		ticket.PhoneNumber = user.PhoneNumber
	}
	created, err := s.Repo.Create(ticket)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *TravelTicketService) GetByID(id int64) (*tentity.TravelTicket, error) {
	return s.Repo.GetByID(id)
}

func (s *TravelTicketService) GetAll() ([]tentity.TravelTicket, error) {
	return s.Repo.GetAll()
}

func (s *TravelTicketService) Update(id int64, dto *models.TravelTicketUpdateDto) (*tentity.TravelTicket, error) {
	ticket, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	ticket = mapper.ApplyUpdateDtoToEntity(dto, ticket)
	updated, err := s.Repo.Update(ticket)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *TravelTicketService) Delete(id int64) error {
	return s.Repo.Delete(id)
}

func (s *TravelTicketService) GetUserResponse(id int64) (*models.TravelTicketUserResponseDto, error) {
	ticket, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	user, err := s.UserRepo.GetByID(ticket.UserID)
	if err != nil {
		return nil, err
	}
	return mapper.ToUserResponseDto(ticket, user), nil
}

func (s *TravelTicketService) GetUserResponses(userID int64) ([]*models.TravelTicketUserResponseDto, error) {
	tickets, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	responses := make([]*models.TravelTicketUserResponseDto, 0, len(tickets))
	for i := range tickets {
		responses = append(responses, mapper.ToUserResponseDto(&tickets[i], user))
	}
	return responses, nil
}

// RecommendForTicket computes best match, best group, and other alternatives
func (s *TravelTicketService) RecommendForTicket(ticketID int64) (*models.RecommendationResult, error) {
	t, err := s.Repo.GetByID(ticketID)
	if err != nil {
		return nil, err
	}

	// Use ticket's local day for same-date matching
	loc := t.DepartureAt.Location()
	day := time.Date(t.DepartureAt.Year(), t.DepartureAt.Month(), t.DepartureAt.Day(), 0, 0, 0, 0, loc)

	var candidates []tentity.TravelTicket
	if models.IsHostel(t.Destination) {
		// Return trip: Home → Hostel
		candidates, err = s.Repo.GetCandidatesSameDateReturn(t.Source, day, t.ID)
		if err != nil {
			return nil, err
		}
	} else {
		// Outbound trip: Hostel → Home
		candidates, err = s.Repo.GetCandidatesSameDateOutbound(t.Destination, day, t.ID)
		if err != nil {
			return nil, err
		}
	}

	// Filter out same user tickets
	filteredCandidates := make([]tentity.TravelTicket, 0, len(candidates))
	for _, c := range candidates {
		if c.UserID != t.UserID {
			filteredCandidates = append(filteredCandidates, c)
		}
	}
	candidates = filteredCandidates

	// Score all candidates
	scored := make([]models.ScoredTicket, 0, len(candidates))
	for _, c := range candidates {
		score := s.scoreTicket(*t, c)
		scored = append(scored, models.ScoredTicket{
			Ticket: c,
			Score:  score,
			Date:   c.DepartureAt.Format("2006-01-02"),
			Time:   c.DepartureAt.Format("15:04"),
		})
	}

	sort.Slice(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })

	result := &models.RecommendationResult{}
	if len(scored) > 0 {
		result.BestMatch = &scored[0]
	}

	// Build Best Group: greedy, 2-hour window around anchor, max 6 users
	group := make([]models.ScoredTicket, 0, 4)
	for _, sct := range scored {
		if len(group) >= 4 {
			break
		}
		c := findCandidateByID(candidates, sct.Ticket.ID)
		if c == nil {
			continue
		}
		if absDuration(c.DepartureAt.Sub(t.DepartureAt)) <= 2*time.Hour {
			group = append(group, sct)
		}
	}
	if len(group) >= 2 {
		result.BestGroup = group
	}

	// Other alternatives
	others := make([]models.ScoredTicket, 0)
	for _, sct := range scored {
		if containsTicketID(result.BestGroup, sct.Ticket.ID) {
			continue
		}
		if result.BestMatch != nil && result.BestMatch.Ticket.ID == sct.Ticket.ID {
			continue
		}
		others = append(others, sct)
	}
	result.OtherAlternatives = others

	return result, nil
}

// helper scoring and filters
func (s *TravelTicketService) scoreTicket(target, candidate tentity.TravelTicket) float64 {
	score := 100.0

	// Time difference penalty
	diffMins := math.Abs(candidate.DepartureAt.Sub(target.DepartureAt).Minutes())
	score -= 0.5 * diffMins
	if score < 0 {
		score = 0
	}

	// Source / destination weighting
	if models.IsHostel(target.Destination) {
		// Return: Home → Hostel
		if target.Source != candidate.Source {
			if models.IsAirportTerminal(target.Source) && models.IsAirportTerminal(candidate.Source) && models.AreNearbyTerminals(target.Source, candidate.Source) {
				score *= 0.8 // nearby terminal less penalty
			} else {
				score *= 0.0 // different source, no match
			}
		}
		if target.Destination != candidate.Destination {
			if models.AreNearbyHostels(target.Destination, candidate.Destination) {
				score *= 0.7 // nearby hostel
			} else {
				score -= 20 // different hostel
			}
		}
	} else {
		// Outbound: Hostel → Home
		if target.Source != candidate.Source {
			if models.AreNearbyHostels(target.Source, candidate.Source) {
				score *= 0.85 // nearby hostel less priority
			} else {
				score *= 0.0 // different hostel, invalid
			}
		}
		if target.Destination != candidate.Destination {
			if models.AreNearbyTerminals(target.Destination, candidate.Destination) {
				score *= 0.6 // nearby terminal less priority - more penalty
			} else {
				score -= 20
			}
		}
	}

	if score < 0 {
		score = 0
	}
	return score
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func containsTicketID(list []models.ScoredTicket, id int64) bool {
	for _, s := range list {
		if s.Ticket.ID == id {
			return true
		}
	}
	return false
}

func findCandidateByID(cands []tentity.TravelTicket, id int64) *tentity.TravelTicket {
	for i := range cands {
		if cands[i].ID == id {
			return &cands[i]
		}
	}
	return nil
}
