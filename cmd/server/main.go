/*
Package main is the entry point of the Go User API application.

It is responsible for:
1. Loading configuration from environment variables.
2. Initializing the structured logger (Zap).
3. Establishing a connection to the PostgreSQL database.
4. setting up the dependency injection container (Repository -> Service -> Handler).
5. Configuring the GoFiber HTTP server and middleware.
6. Registering API routes and starting the server.
*/
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rohanparmar/go-user-api/config"
	db "github.com/rohanparmar/go-user-api/db/sqlc/generated"
	"github.com/rohanparmar/go-user-api/internal/handler"
	"github.com/rohanparmar/go-user-api/internal/logger"
	"github.com/rohanparmar/go-user-api/internal/middleware"
	"github.com/rohanparmar/go-user-api/internal/repository"
	"github.com/rohanparmar/go-user-api/internal/routes"
	"github.com/rohanparmar/go-user-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize logger
	env := cfg.GetEnv("ENV", "development")
	if err := logger.InitLogger(env); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting Go User API server...")

	// Connect to PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	logger.Log.Info("Database connection established successfully")

	// Initialize SQLC queries
	queries := db.New(pool)

	// Initialize layers (Repository -> Service -> Handler)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Log.Error("Request error", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	// Middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestDuration())

	// Setup routes
	routes.SetupRoutes(app, userHandler)

	// Start server
	port := cfg.GetEnv("PORT", "8080")
	logger.Log.Info("Server starting", zap.String("port", port))
	
	if err := app.Listen(":" + port); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}

