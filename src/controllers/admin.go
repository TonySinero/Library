
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/models"
	"net/http"
)


type adminResponse struct {
	Message string `json:"message"`
}


// Create admin.

func CreateAdmin(c *gin.Context) {
	// Validate input
	var input models.AdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin := models.Admin{
		Email:    input.Email,
		Password: input.Password,
	}
	admin.EncryptPassword()
	models.DB.Create(&admin)

	c.JSON(http.StatusCreated, adminResponse{"OK"})
}
