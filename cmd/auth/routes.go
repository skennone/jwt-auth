package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	router.POST("/api/login", app.Login)
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

func (app *application) Login(c *gin.Context) {
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
	user, err := app.models.Users.GetByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	match, err := user.Password.Matches(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if !match {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}
