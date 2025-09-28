package hander

import (
	"Travel_Sync/internal/user/models"
	"Travel_Sync/internal/user/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (u *UserHandler) CreateUser(c *gin.Context) {
	var userCreateDto models.UserCreateDto
	if err := c.ShouldBind(&userCreateDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Failed to bind request body for user create method"})
		return
	}

	user, err := u.svc.CreateUser(&userCreateDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "user": user})
}

func (u *UserHandler) UpdateUser(c *gin.Context) {
	var userUpdateDto models.UserUpdateDto
	if err := c.ShouldBind(&userUpdateDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Failed to bind request body for User update method"})
		return
	}
	err := u.svc.UpdateUser()
}
