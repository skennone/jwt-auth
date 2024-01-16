package main

import (
	"github.com/gin-gonic/gin"
	"github.com/skennone/goAuth/controllers"
)

type Config struct {
	R *gin.Engine
}
type Handler struct{}

func routes() {
	router := gin.Default()

	router.POST("/api/register", controllers.Register)
	router.Run(":8080")
}
