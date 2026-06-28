package handler

import (
	"net/http"
	"rate-limiter/config"
	"rate-limiter/limiter"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type limitConfig struct {
	Path              string
	TargetService     string
	LimitStrategyName string
	Limiter           limiter.RequestLimiter
	Config            config.StrategyConfig
}

type Limiter struct {
	LimitRegistery map[string]limitConfig
	Logger         *zap.Logger
}

func NewLimiter(l *zap.Logger) *Limiter {
	return &Limiter{
		LimitRegistery: make(map[string]limitConfig),
		Logger:         l,
	}
}

type EndpointRegisterReq struct {
	Path          string                `json:"path"`
	Method        string                `json:"method"`
	Strategy      config.StrategyConfig `json:"strategy"`
	TargetService string                `json:"target_service"`
}

func (l *Limiter) RegisterEndpoint(c *gin.Context) {
	var req EndpointRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rateLimiterKey := req.Method + ":" + req.Path

	if _, found := l.LimitRegistery[rateLimiterKey]; found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint already registered"})
		return
	}

	limitStrategyName := req.Strategy.Name
	reqLimiter := limiter.GetRequestLimiter(limitStrategyName)

	if reqLimiter == nil {
		l.Logger.Error("Invalid strategy name", zap.String("strategy", limitStrategyName))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy name"})
		return
	}

	strategyConfig := config.StrategyConfig{
		Name:     req.Strategy.Name,
		Limit:    req.Strategy.Limit,
		Interval: req.Strategy.Interval,
	}

	l.LimitRegistery[rateLimiterKey] = limitConfig{
		Path:              req.Path,
		LimitStrategyName: req.Strategy.Name,
		Limiter:           reqLimiter,
		Config:            strategyConfig,
		TargetService:     req.TargetService,
	}

	l.Logger.Info("Endpoint registered successfully", zap.String("endpoint", rateLimiterKey))

	// TODO: Add endpoint to database
	// TODO: Add strategy to cache

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
