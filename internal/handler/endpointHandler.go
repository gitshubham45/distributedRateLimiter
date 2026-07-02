package handler

import (
	"context"
	"net/http"
	"rate-limiter/internal/config"
	"rate-limiter/internal/domain"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/repository"

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
	EndpointRepo   repository.EndpointRepository
}

func NewLimiter(l *zap.Logger, endpointRepo repository.EndpointRepository) *Limiter {
	return &Limiter{
		LimitRegistery: make(map[string]limitConfig),
		Logger:         l,
		EndpointRepo:   endpointRepo,
	}
}

func (l *Limiter) LoadRegisteredEndpoints(ctx context.Context) error {
	endpoints, err := l.EndpointRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		reqLimiter := limiter.GetRequestLimiter(endpoint.StrategyName)
		strategyConfig := config.StrategyConfig{
			Name:     endpoint.StrategyName,
			Limit:    endpoint.Limit,
			Interval: endpoint.Interval,
		}

		l.LimitRegistery[endpoint.Key] = limitConfig{
			Path:              endpoint.Path,
			TargetService:     endpoint.TargetService,
			LimitStrategyName: endpoint.StrategyName,
			Limiter:           reqLimiter,
			Config:            strategyConfig,
		}
	}

	l.Logger.Info("Registered endpoints loaded", zap.Int("count", len(endpoints)), zap.Any("enspoints : ", endpoints))
	return nil
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

	endpointConfig := domain.EndpointConfig{
		Key:           rateLimiterKey,
		Method:        req.Method,
		Path:          req.Path,
		TargetService: req.TargetService,
		StrategyName:  req.Strategy.Name,
		Limit:         req.Strategy.Limit,
		Interval:      req.Strategy.Interval,
	}

	if err := l.EndpointRepo.Save(c.Request.Context(), endpointConfig); err != nil {
		l.Logger.Error("Failed to save endpoint", zap.String("endpoint", rateLimiterKey), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save endpoint"})
		return
	}

	l.LimitRegistery[rateLimiterKey] = limitConfig{
		Path:              req.Path,
		LimitStrategyName: req.Strategy.Name,
		Limiter:           reqLimiter,
		Config:            strategyConfig,
		TargetService:     req.TargetService,
	}

	l.Logger.Info("Endpoint registered successfully", zap.String("endpoint", rateLimiterKey))

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
