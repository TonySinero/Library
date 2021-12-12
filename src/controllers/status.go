
package controllers

import (
	"github.com/gin-gonic/gin"
)

type statusMessage struct {
	Message string `json:"message"`
}


//  Health api.

func Status(c *gin.Context) {
	c.JSON(
		200,
		statusMessage{"Successfully"},
	)
}
