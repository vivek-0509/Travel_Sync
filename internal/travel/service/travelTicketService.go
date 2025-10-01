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
	// Ensure user has no other ticket on the same date
	day := time.Date(ticket.DepartureAt.Year(), ticket.DepartureAt.Month(), ticket.DepartureAt.Day(), 0, 0, 0, 0, ticket.DepartureAt.Location())
	exists, err := s.Repo.ExistsForUserOnDate(userID, day, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ticket already exists for this date")
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

func (s *TravelTicketService) Update(currentUserID int64, id int64, dto *models.TravelTicketUpdateDto) (*tentity.TravelTicket, error) {
	ticket, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ticket.UserID != currentUserID {
		return nil, errors.New("forbidden")
	}
	ticket = mapper.ApplyUpdateDtoToEntity(dto, ticket)
	// If departure time changed (or even if not), enforce single ticket per date
	day := time.Date(ticket.DepartureAt.Year(), ticket.DepartureAt.Month(), ticket.DepartureAt.Day(), 0, 0, 0, 0, ticket.DepartureAt.Location())
	excludeID := id
	exists, err := s.Repo.ExistsForUserOnDate(currentUserID, day, &excludeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ticket already exists for this date")
	}
	updated, err := s.Repo.Update(ticket)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *TravelTicketService) Delete(currentUserID int64, id int64) error {
	ticket, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if ticket.UserID != currentUserID {
		return errors.New("forbidden")
	}
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

// GetByUser returns all tickets created by the specified user
func (s *TravelTicketService) GetByUser(userID int64) ([]tentity.TravelTicket, error) {
	return s.Repo.GetByUserID(userID)
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

	// Score all candidates, but only within the asymmetric time window
	beforeWindow := time.Duration(t.TimeDiffMins) * time.Minute
	afterWindow := 60 * time.Minute
	scored := make([]models.ScoredTicket, 0, len(candidates))
	for _, c := range candidates {
		delta := c.DepartureAt.Sub(t.DepartureAt)
		if delta <= 0 {
			// Candidate leaves before or at the same time as target
			if absDuration(delta) > beforeWindow {
				continue
			}
		} else {
			// Candidate leaves after target
			if delta > afterWindow {
				continue
			}
		}

		score := s.scoreTicket(*t, c)
		// fetch minimal user details for candidate
		cu, uerr := s.UserRepo.GetByID(c.UserID)
		var minUser models.MinimalUser
		if uerr == nil && cu != nil {
			minUser = models.MinimalUser{Name: cu.Name, Batch: cu.Batch, Email: cu.Email}
		}
		public := models.PublicTicket{
			Source:       c.Source,
			Destination:  c.Destination,
			EmptySeats:   c.EmptySeats,
			DepartureAt:  c.DepartureAt,
			TimeDiffMins: c.TimeDiffMins,
			PhoneNumber:  c.PhoneNumber,
			Status:       c.Status,
			CreatedAt:    c.CreatedAt,
			UpdatedAt:    c.UpdatedAt,
		}
		scored = append(scored, models.ScoredTicket{
			Ticket:      public,
			Score:       score,
			Date:        c.DepartureAt.Format("2006-01-02"),
			Time:        c.DepartureAt.Format("15:04"),
			User:        minUser,
			CandidateID: c.ID,
		})
	}

	sort.Slice(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })

	result := &models.RecommendationResult{}
	if len(scored) > 0 {
		result.BestMatch = &scored[0]
	}

	// Build Best Group: greedy, asymmetric window (before: TimeDiffMins; after: 60m)
	group := make([]models.ScoredTicket, 0, 4)
	timeWindowBefore := beforeWindow
	timeWindowAfter := afterWindow
	for _, sct := range scored {
		if len(group) >= 4 {
			break
		}
		c := findCandidateByID(candidates, sct.CandidateID)
		if c == nil {
			continue
		}
		delta := c.DepartureAt.Sub(t.DepartureAt)
		if delta <= 0 {
			if absDuration(delta) <= timeWindowBefore {
				group = append(group, sct)
			}
		} else {
			if delta <= timeWindowAfter {
				group = append(group, sct)
			}
		}
	}
	if len(group) >= 2 {
		result.BestGroup = group
	}

	// Other alternatives
	others := make([]models.ScoredTicket, 0)
	for _, sct := range scored {
		if containsTicketID(result.BestGroup, sct.CandidateID) {
			continue
		}
		if result.BestMatch != nil && result.BestMatch.CandidateID == sct.CandidateID {
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
		if s.CandidateID == id {
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
