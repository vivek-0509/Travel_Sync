package mapper

import (
	tentity "Travel_Sync/internal/travel/entity"
	"Travel_Sync/internal/travel/models"
	uentity "Travel_Sync/internal/user/entity"
	"time"
)

func FromCreateDtoToEntity(dto *models.TravelTicketCreateDto, user *uentity.User) (*tentity.TravelTicket, error) {
	departureAt, err := time.Parse(time.RFC3339, dto.DepartureAt)
	if err != nil {
		return nil, err
	}
	return &tentity.TravelTicket{
		Source:       dto.Source,
		Destination:  dto.Destination,
		EmptySeats:   dto.EmptySeats,
		DepartureAt:  departureAt,
		TimeDiffMins: dto.TimeDiffMins,
		UserID:       user.ID,
		PhoneNumber:  dto.PhoneNumber,
		Status:       "open",
	}, nil
}

func ApplyUpdateDtoToEntity(dto *models.TravelTicketUpdateDto, ticket *tentity.TravelTicket) *tentity.TravelTicket {
	if dto.Source != "" {
		ticket.Source = dto.Source
	}
	if dto.Destination != "" {
		ticket.Destination = dto.Destination
	}
	if dto.DepartureAt != "" {
		if t, err := time.Parse(time.RFC3339, dto.DepartureAt); err == nil {
			ticket.DepartureAt = t
		}
	}
	if dto.TimeDiffMins != 0 {
		ticket.TimeDiffMins = dto.TimeDiffMins
	}
	if dto.EmptySeats != 0 {
		ticket.EmptySeats = dto.EmptySeats
	}
	if dto.PhoneNumber != "" {
		ticket.PhoneNumber = dto.PhoneNumber
	}
	if dto.Status != "" {
		ticket.Status = dto.Status
	}
	return ticket
}

func ToUserResponseDto(ticket *tentity.TravelTicket, user *uentity.User) *models.TravelTicketUserResponseDto {
	return &models.TravelTicketUserResponseDto{
		ID:           ticket.ID,
		StudentName:  user.Name,
		StudentBatch: user.Batch,
		Source:       ticket.Source,
		Destination:  ticket.Destination,
		Date:         ticket.DepartureAt.Format("2006-01-02"),
		Time:         ticket.DepartureAt.Format("15:04"),
		EmptySeats:   ticket.EmptySeats,
		PhoneNumber:  ticket.PhoneNumber,
	}
}

// no defaultStatus helper needed once create always sets to "open"
