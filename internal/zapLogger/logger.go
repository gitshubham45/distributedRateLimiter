package zapLogger

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() *zap.Logger {

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err.Error())
	}

	encodingConfig := zap.NewProductionEncoderConfig()
	encodingConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encodingConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(logFile), zap.DebugLevel),
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	logger = zap.New(core)

	return logger
}

func ZapLoggerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Execute the next handlers down the line
		c.Next()

		// Calculate latency after execution completes
		latency := time.Since(start)

		// Output structural logs in JSON format
		logger.Info("HTTP REQUEST",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			// Latency is explicitly serialized to seconds in production configs
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}
