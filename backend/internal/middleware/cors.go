package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "X-Request-ID"}
	return cors.New(config)
}
