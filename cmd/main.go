package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shuter-go/config"
	"shuter-go/internal/handlers"
	"shuter-go/internal/routes"
	"shuter-go/internal/services"
	"shuter-go/internal/storage"
	"shuter-go/pkg/logger"
	"syscall"
	"time"
)

var cfg = config.NewServerConfig()

func main() {
	if err := logger.Init("INFO"); err != nil {
		panic(err)
	}
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	//authMW := middlewares.AuthMiddleware([]byte(jwtSecret))

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	db, err := storage.NewDB(storage.Config{
		DBCfg:      cfg.DBCfg,
		Migrations: cfg.Migrations,
	})
	if err != nil {
		logger.Log.Error("Failed to initialize database", zap.Error(err))
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Log.Info("Error DB close", zap.String("address", cfg.RunAddr))
		}
	}(db)

	r := gin.Default()

	userService := services.New()
	userHandler := handlers.New(userService)

	routes.SetupUserRoutes(r, userHandler)

	srv := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: r,
	}

	serverErr := make(chan error, 1)

	go func() {
		logger.Log.Info("Server is running", zap.String("address", cfg.RunAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Log.Info("Shutting down server...")
	case err := <-serverErr:
		logger.Log.Error("Server error", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server is shutdown", zap.Error(err))
	}

	logger.Log.Info("Server stopped gracefully")
}
