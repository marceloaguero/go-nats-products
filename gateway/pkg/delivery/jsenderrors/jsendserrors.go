package jsenderrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": message,
	})
}
