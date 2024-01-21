package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skennone/goAuth/internal/data"
	"github.com/skennone/goAuth/internal/validator"
)

type Config struct {
	R *gin.Engine
}
type Handler struct{}

func (app *application) routes() {
	router := gin.Default()

	router.POST("/api/register", app.Register)
	router.Run(":8080")
}

func (app *application) Register(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user := &data.User{
		Email: body.Email,
	}

	err := user.Password.Set(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Server error",
		})
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Server error",
		})
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Server error",
		})
	}
}
