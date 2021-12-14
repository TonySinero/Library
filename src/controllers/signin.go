
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library/src/models"
	"github.com/library/src/services"
	"net/http"
)

type signInReponse struct {
	Token string `json:"token"`
}

type signInModel struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (lm *signInModel) Validate() bool {
	var admin models.Admin
	if err := models.DB.First(&admin, "email = ?", lm.Email).Error; err != nil {
		return false
	}
	return admin.ValidatePassword(lm.Password)
}

// SignIn Controller.

func SignIn(c *gin.Context) {
	serviceJWT := services.JWTAuthService()

	var input signInModel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !input.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password invalid."})
		return
	}
	token := serviceJWT.GenerateToken(input.Email, true)
	c.JSON(http.StatusOK, signInReponse{token})
}
