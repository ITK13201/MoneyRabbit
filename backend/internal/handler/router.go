package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itk13201/money-rabbit/internal/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health", health)
	}

	return r
}

// @Summary      ヘルスチェック
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
