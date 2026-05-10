package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/itk13201/money-rabbit/internal/middleware"
)

// Deps holds all handler dependencies wired from main.
type Deps struct {
	Category    *CategoryHandler
	Import      *ImportHandler
	Transaction *TransactionHandler
	Summary     *SummaryHandler
}

func NewRouter(deps Deps) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Swagger UI
	r.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/health", health)

		// Import formats (read-only, no DB)
		api.GET("/import-formats", ListFormats)

		// CSV import
		api.POST("/import/preview", Preview)
		api.POST("/import/confirm", deps.Import.Confirm)

		// Categories
		cats := api.Group("/categories")
		{
			cats.GET("", deps.Category.List)
			cats.POST("", deps.Category.Create)
			cats.GET("/:id", deps.Category.Get)
			cats.PUT("/:id", deps.Category.Update)
			cats.DELETE("/:id", deps.Category.Delete)
		}

		// Category rules
		rules := api.Group("/category-rules")
		{
			rules.POST("", deps.Category.CreateRule)
			rules.PUT("/:id", deps.Category.UpdateRule)
			rules.DELETE("/:id", deps.Category.DeleteRule)
		}

		// Transactions
		txs := api.Group("/transactions")
		{
			txs.GET("", deps.Transaction.List)
			txs.PATCH("/:id/category", deps.Transaction.UpdateCategory)
			txs.DELETE("/:id", deps.Transaction.Delete)
		}

		// Summary
		summary := api.Group("/summary")
		{
			summary.GET("/monthly", deps.Summary.Monthly)
		}
	}

	return r
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
