package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/models"
	"net/http"
)


// Get all books

func FindBooks(c *gin.Context) {
	var books []models.Book
	models.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}


// Description Detail of a book by ID

func DetailBooks(c *gin.Context) {
	var book models.Book

	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found!"})
		return
	}

	c.JSON(http.StatusOK, book)
}


// Create a book.

func CreateBook(c *gin.Context) {
	// Validate input
	var input models.CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := models.Book{
		Title:          input.Title,
		Price:          input.Price,
		Author:         input.Author,
		Year:           input.Year,
		NumberOfCopies: input.NumberOfCopies,
		NumberOfPages:  input.NumberOfPages,
		Image:          input.Image,
		PricePerDay:    input.PricePerDay,
		CategoryID:     input.CategoryID,
	}
	models.DB.Create(&book)

	c.JSON(http.StatusCreated, input)
}

// Delete a book.

func DeleteBook(c *gin.Context) {
	// Get model if exist
	var books models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&books).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "books not found!"})
		return
	}
	models.DB.Delete(&books)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// Update a book.

func UpdateBook(c *gin.Context) {
	// Get model if exist
	var books models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&books).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "books not found!"})
		return
	}
	// Validate input
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&books).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": books})
}
