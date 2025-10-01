package handler

import (
	"net/http"
	"strconv"

	"Travel_Sync/internal/travel/models"
	tservice "Travel_Sync/internal/travel/service"

	"github.com/gin-gonic/gin"
)

type TravelTicketHandler struct {
	Svc *tservice.TravelTicketService
}

func NewTravelTicketHandler(svc *tservice.TravelTicketService) *TravelTicketHandler {
	return &TravelTicketHandler{Svc: svc}
}

func parseID(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id is required"})
		return 0, false
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return 0, false
	}
	return id, true
}

// routes are now registered in routes package

func (h *TravelTicketHandler) Create(c *gin.Context) {
	claims, _ := c.Get("jwt_claims")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}
	_ = claims

	var dto models.TravelTicketCreateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request body"})
		return
	}

	ticket, err := h.Svc.Create(userID.(int64), &dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": ticket})
}

func (h *TravelTicketHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ticket, err := h.Svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "ticket not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ticket})
}

func (h *TravelTicketHandler) GetRecommendations(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	result, err := h.Svc.RecommendForTicket(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *TravelTicketHandler) GetAll(c *gin.Context) {
	tickets, err := h.Svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to fetch tickets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tickets})
}

// GetMyTickets returns all tickets created by the authenticated user
func (h *TravelTicketHandler) GetMyTickets(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}
	var userID int64
	switch v := uid.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid user context"})
		return
	}

	tickets, err := h.Svc.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tickets})
}

func (h *TravelTicketHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var dto models.TravelTicketUpdateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request body"})
		return
	}
	ticket, err := h.Svc.Update(id, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ticket})
}

func (h *TravelTicketHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.Svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to delete ticket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": "ticket deleted"})
}

func (h *TravelTicketHandler) GetUserResponses(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}
	userID, ok := uid.(int64)
	if !ok {
		// JWT middleware may store numeric as float64; handle conversion
		if f, ok2 := uid.(float64); ok2 {
			userID = int64(f)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid user context"})
			return
		}
	}

	responses, err := h.Svc.GetUserResponses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": responses})
}
