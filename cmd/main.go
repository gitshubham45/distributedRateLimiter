package main

import (
	"context"
	"os"
	"rate-limiter/internal/config"
	"rate-limiter/internal/database"
	"rate-limiter/internal/handler"
	"rate-limiter/internal/repository"
	"rate-limiter/internal/zapLogger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
1. Rate limiter (not distributed)

// request -> check limit -> if allowable or not


*/

func main() {

	logger := zapLogger.InitLogger()

	defer logger.Sync()

	if err := config.InitConfig(configPath()); err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	// Connect to database
	db := database.NewDB(config.DbConfig)
	endpointRepo := repository.NewPostgresEndpointRepository(db)
	if err := endpointRepo.AutoMigrate(); err != nil {
		logger.Fatal("error migrating endpoint table", zap.Error(err))
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(zapLogger.ZapLoggerMiddleWare())
	r.Use(gin.Recovery())

	limiter := handler.NewLimiter(logger, endpointRepo)
	if err := limiter.LoadRegisteredEndpoints(context.Background()); err != nil {
		logger.Fatal("error loading registered endpoints", zap.Error(err))
	}

	r.Any("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.POST("/register-endpoint", limiter.RegisterEndpoint)
	r.NoRoute(limiter.HandleLimit)

	port := ":8080"

	logger.Info("starting server", zap.String("port", port))
	err := r.Run(port)
	if err != nil {
		logger.Fatal("error starting server", zap.Error(err))
	}
	logger.Info("server started at", zap.String("port", port))
}

func configPath() string {
	if path := os.Getenv("RATE_LIMITER_CONFIG_PATH"); path != "" {
		return path
	}

	if _, err := os.Stat("rate-limiter.yaml"); err == nil {
		return "rate-limiter.yaml"
	}

	return "../rate-limiter.yaml"
}
