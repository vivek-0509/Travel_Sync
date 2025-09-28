package hander

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
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id is required"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid id"})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id cannot be 0"})
		return
	}

	var userUpdateDto models.UserUpdateDto
	if err := c.ShouldBind(&userUpdateDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Failed to bind request body for User update method"})
		return
	}

	user, err := u.svc.UpdateUser(id, &userUpdateDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func (u *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id is required"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid id"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id cannot be 0"})
		return
	}

	err = u.svc.DeleteByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to delete user"})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": "User deleted successfully"})
}

func (u *UserHandler) GetUserById(c *gin.Context) {
	idstr := c.Param("id")
	if idstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id is required"})
		return
	}
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid id"})
		return
	}

	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id cannot be 0"})
		return
	}

	user, err := u.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})

}

func (u *UserHandler) GetAllUser(c *gin.Context) {
	users, err := u.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get all users"})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}
