package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/models"
	"net/http"
)


// Get all users

func GetUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}


// Description Detail of a user by ID

func DetailUser(c *gin.Context) {
	var user models.User

	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	c.JSON(http.StatusOK, user)
}


// Create a users.

func CreateUsers(c *gin.Context) {
	// Validate input
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := models.User{
		Name:         input.Name,
		Surname:      input.Surname,
		SecondName:   input.SecondName,
		Passport:     input.Passport,
		DateOfBirth:  input.DateOfBirth,
		Email:        input.Email,
		Address:      input.Address,
		Indebtedness: input.Indebtedness,

	}
	models.DB.Create(&user)

	c.JSON(http.StatusCreated, input)
}


// Update a user.

func UpdateUser(c *gin.Context) {
	// Get model if exist
	var user models.User
	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found!"})
		return
	}
	// Validate input
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&user).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": user})
}