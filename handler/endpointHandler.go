package handler

import (
	"net/http"
	"rate-limiter/limiter"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type limitConfigStore struct {
	path              string
	limitStrategyName string
	limiter           limiter.RequestLimiter
	limitConfig       StrategyConfig
}

type Limiter struct {
	limitRegistery map[string]limitConfigStore
	logger         *zap.Logger
}

func NewLimiter(l *zap.Logger) *Limiter {
	return &Limiter{
		limitRegistery: make(map[string]limitConfigStore),
		logger:         l,
	}
}

type StrategyConfig struct {
	Name     string `json:"name"`
	Limit    int64  `json:"limit"`
	Interval int64  `json:"interval"`
}

type EndpointRegisterReq struct {
	Path     string         `json:"path"`
	Method   string         `json:"method"`
	Strategy StrategyConfig `json:"strategy"`
}

func (l *Limiter) RegisterEndpoint(c *gin.Context) {
	var req EndpointRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limitKey := req.Method + ":" + req.Path

	if _, found := l.limitRegistery[limitKey]; found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint already registered"})
		return
	}

	limitStrategyName := req.Strategy.Name
	reqLimiter := limiter.GetRequestLimiter(limitStrategyName)

	if reqLimiter == nil {
		l.logger.Error("Invalid strategy name", zap.String("strategy", limitStrategyName))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy name"})
		return
	}

	strategyConfig := StrategyConfig{
		Name:     req.Strategy.Name,
		Limit:    req.Strategy.Limit,
		Interval: req.Strategy.Interval,
	}

	l.limitRegistery[limitKey] = limitConfigStore{
		path:              req.Path,
		limitStrategyName: req.Strategy.Name,
		limiter:           reqLimiter,
		limitConfig:       strategyConfig,
	}

	l.logger.Info("Endpoint registered successfully", zap.String("endpoint", limitKey))

	// TODO: Add endpoint to database
	// TODO: Add strategy to cache

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
