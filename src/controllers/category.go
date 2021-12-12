
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/models"
	"net/http"
)


// Create a category.

func CreateCategory(c *gin.Context) {
	// Validate input
	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create category
	category := models.Category{Name: input.Name}
	models.DB.Create(&category)

	c.JSON(http.StatusCreated, input)
}


//  List categories.

func FindCategory(c *gin.Context) {
	var categories []models.Category
	models.DB.Find(&categories)

	c.JSON(http.StatusOK, categories)
}
