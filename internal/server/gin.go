package server

import "github.com/gin-gonic/gin"

func NewGinRouter() *gin.Engine {
	r := gin.Default()
	return r
}
