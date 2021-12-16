package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/database"
	"github.com/library/src/models"
	"net/http"
)

func GetAllAccept(c *gin.Context) {
	var accept []models.Acceptance
	database.DB.Find(&accept)

	c.JSON(http.StatusOK, gin.H{"data": accept})
}


func GetAccept(c *gin.Context) {
	var form models.Acceptance

	if err := database.DB.Where("id = ?", c.Param("id")).First(&form).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form not found!"})
		return
	}

	c.JSON(http.StatusOK, form)
}


// Create a form.

func CreateAccept(c *gin.Context) {
	// Validate input
	var input models.CreateAcceptanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create form
	form := models.Acceptance{
		Name:         input.Name,
		Surname:      input.Surname,
		SecondName:   input.SecondName,
		Passport:     input.Passport,
		Books:        input.Books,
		Condition:    input.Condition,
		Rating:       input.Rating,
		FinalPrice:   input.FinalPrice,
	}
	database.DB.Create(&form)

	c.JSON(http.StatusCreated, input)
}

// Delete a form.

func DeleteAccept(c *gin.Context) {
	// Get model if exist
	var form models.Acceptance
	if err := database.DB.Where("id = ?", c.Param("id")).First(&form).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form not found!"})
		return
	}
	database.DB.Delete(&form)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
