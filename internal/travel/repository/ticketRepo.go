package repository

import (
	"Travel_Sync/internal/travel/entity"
	"Travel_Sync/internal/travel/models"
	"time"

	"gorm.io/gorm"
)

type TravelTicketRepo struct {
	DB *gorm.DB
}

func NewTravelTicketRepo(db *gorm.DB) *TravelTicketRepo {
	return &TravelTicketRepo{DB: db}
}

func (r *TravelTicketRepo) Create(ticket *entity.TravelTicket) (*entity.TravelTicket, error) {
	if err := r.DB.Create(ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TravelTicketRepo) GetByID(id int64) (*entity.TravelTicket, error) {
	var ticket entity.TravelTicket
	if err := r.DB.First(&ticket, id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TravelTicketRepo) GetAll() ([]entity.TravelTicket, error) {
	var tickets []entity.TravelTicket
	if err := r.DB.Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TravelTicketRepo) Update(ticket *entity.TravelTicket) (*entity.TravelTicket, error) {
	if err := r.DB.Save(ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TravelTicketRepo) Delete(id int64) error {
	return r.DB.Delete(&entity.TravelTicket{ID: id}).Error
}

func (r *TravelTicketRepo) GetByUserID(userID int64) ([]entity.TravelTicket, error) {
	var tickets []entity.TravelTicket
	if err := r.DB.Where("user_id = ?", userID).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TravelTicketRepo) GetCandidatesSameDateOutbound(destination string, dayStart time.Time, excludeID int64) ([]entity.TravelTicket, error) {
	var tickets []entity.TravelTicket
	dayEnd := dayStart.Add(24 * time.Hour)
	q := r.DB
	if models.IsAirportTerminal(destination) {
		terminals := make([]string, 0, len(models.AirportTerminals))
		for t := range models.AirportTerminals {
			terminals = append(terminals, t)
		}
		q = q.Where("destination IN ?", terminals)
	} else {
		q = q.Where("destination = ?", destination)
	}
	err := q.Where("departure_at >= ? AND departure_at < ? AND id <> ?",
		dayStart, dayEnd, excludeID).Find(&tickets).Error
	return tickets, err
}

func (r *TravelTicketRepo) GetCandidatesSameDateReturn(source string, dayStart time.Time, excludeID int64) ([]entity.TravelTicket, error) {
	var tickets []entity.TravelTicket
	dayEnd := dayStart.Add(24 * time.Hour)
	q := r.DB
	if models.IsAirportTerminal(source) {
		terminals := make([]string, 0, len(models.AirportTerminals))
		for t := range models.AirportTerminals {
			terminals = append(terminals, t)
		}
		q = q.Where("source IN ?", terminals)
	} else {
		q = q.Where("source = ?", source)
	}
	err := q.Where("departure_at >= ? AND departure_at < ? AND id <> ?",
		dayStart, dayEnd, excludeID).Find(&tickets).Error
	return tickets, err
}
