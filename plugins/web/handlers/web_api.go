package handlers

import "github.com/gin-gonic/gin"

func (h *WebHandler) HandleWebApi(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Web API",
	})
}
