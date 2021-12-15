package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/models"
	"net/http"
)
func GetAllForm(c *gin.Context) {
	var forms []models.Lending
	models.DB.Find(&forms)

	c.JSON(http.StatusOK, gin.H{"data": forms})
}


func GetForm(c *gin.Context) {
	var form models.Lending

	if err := models.DB.Where("id = ?", c.Param("id")).First(&form).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form not found!"})
		return
	}

	c.JSON(http.StatusOK, form)
}


// Create a form.

func CreateForm(c *gin.Context) {
	// Validate input
	var input models.CreateLendingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create form
	form := models.Lending{
		Name:         input.Name,
		Surname:      input.Surname,
		SecondName:   input.SecondName,
		Passport:     input.Passport,
		Books:        input.Books,
		ReturnDate:   input.ReturnDate,
		Price:        input.Price,
	}
	models.DB.Create(&form)

	c.JSON(http.StatusCreated, input)
}

// Delete a form.

func DeleteForm(c *gin.Context) {
	// Get model if exist
	var form models.Lending
	if err := models.DB.Where("id = ?", c.Param("id")).First(&form).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form not found!"})
		return
	}
	models.DB.Delete(&form)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
