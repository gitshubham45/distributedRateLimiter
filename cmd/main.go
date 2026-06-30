package main

import (
	"rate-limiter/internal/config"
	"rate-limiter/internal/handler"
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

	if err := config.InitConfig("rate-limiter.yaml"); err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(zapLogger.ZapLoggerMiddleWare())
	r.Use(gin.Recovery())

	limiter := handler.NewLimiter(logger)

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
