package handler

import (
	"Travel_Sync/internal/user/models"
	"Travel_Sync/internal/user/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func getIDParam(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id is required"})
		return 0, false
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid id"})
		return 0, false
	}
	return id, true
}

//func (u *UserHandler) CreateUser(c *gin.Context) {
//
//	user, err := u.svc.CreateUser(email)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create user"})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{"success": true, "user": user})
//}

func (u *UserHandler) UpdateUser(c *gin.Context) {
	id, ok := getIDParam(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Id parsing failed"})
		return
	}

	// Ownership check: user can update only their own profile
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}
	currentUserID := toInt64(uid)
	if currentUserID != id {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "forbidden"})
		return
	}

	var dto models.UserUpdateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request body"})
		return
	}

	user, err := u.svc.UpdateUser(id, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func toInt64(v interface{}) int64 {
	if id, ok := v.(int64); ok {
		return id
	}
	if f, ok := v.(float64); ok {
		return int64(f)
	}
	return 0
}

func (u *UserHandler) DeleteUser(c *gin.Context) {
	id, ok := getIDParam(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Id parsing failed"})
		return
	}

	if err := u.svc.DeleteByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": "User deleted successfully"})
}

func (u *UserHandler) GetUserById(c *gin.Context) {
	id, ok := getIDParam(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Id parsing failed"})
		return
	}

	user, err := u.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func (u *UserHandler) GetAllUser(c *gin.Context) {
	users, err := u.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get all users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}
