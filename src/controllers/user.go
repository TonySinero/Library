
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/models"
	"net/http"
)


type userResponse struct {
	Message string `json:"message"`
}


// Create some user.

func CreateUser(c *gin.Context) {
	// Validate input
	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}
	user.EncryptPassword()
	models.DB.Create(&user)

	c.JSON(http.StatusCreated, userResponse{"OK"})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}
