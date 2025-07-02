package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"todo_project/common/log"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found, proceeding with defaults")
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiKey = os.Getenv("API_KEY")
		authHeader := c.GetHeader("X-API-KEY")
		if authHeader == "" || authHeader != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization failed",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}