package main

import (
	"net/http"

	"github.com/alireza-akbarzadeh/restful-app/internal/database"
	"github.com/alireza-akbarzadeh/restful-app/internal/messages"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registeredRoutes struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=3"`
}

func (app *application) registerUser(c *gin.Context) {

	var input registeredRoutes
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrFailedToRegisterUser})
		return
	}
	input.Password = string(hashedPassword)

	user := database.User{
		Email:    input.Email,
		Password: input.Password,
		Name:     input.Name,
	}
	_, err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrFailedToRegisterUser, "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
