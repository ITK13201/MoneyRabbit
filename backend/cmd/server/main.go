package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itk13201/money-rabbit/ent"
	"github.com/itk13201/money-rabbit/internal/handler"
	"github.com/itk13201/money-rabbit/internal/service/classifier"
	"github.com/itk13201/money-rabbit/internal/service/persistence"
	categoryUC "github.com/itk13201/money-rabbit/internal/usecase/category"
	txUC "github.com/itk13201/money-rabbit/internal/usecase/transaction"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	client, err := ent.Open("mysql", dsn)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	// Auto-migrate schema (creates tables if not exist)
	if err := client.Schema.Create(context.Background()); err != nil {
		slog.Error("failed to migrate schema", "error", err)
		os.Exit(1)
	}

	// Persistence layer
	categoryRepo := persistence.NewCategoryRepository(client)
	transactionRepo := persistence.NewTransactionRepository(client)

	// Classifier (Claude API); gracefully handles missing API key
	cl := classifier.New(os.Getenv("ANTHROPIC_API_KEY"))

	// Use cases
	categoryUsecase := categoryUC.New(categoryRepo)
	importUsecase := txUC.NewImportUsecase(transactionRepo, categoryRepo, cl)
	listUsecase := txUC.NewListUsecase(transactionRepo)
	updateCatUsecase := txUC.NewUpdateCategoryUsecase(transactionRepo)

	// Handlers
	deps := handler.Deps{
		Category:    handler.NewCategoryHandler(categoryUsecase),
		Import:      handler.NewImportHandler(importUsecase),
		Transaction: handler.NewTransactionHandler(listUsecase, updateCatUsecase),
	}

	r := handler.NewRouter(deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
