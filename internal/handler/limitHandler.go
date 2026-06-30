package handler

import (
	"net/http"
	"rate-limiter/internal/config"
	"rate-limiter/internal/proxy"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (l *Limiter) HandleLimit(c *gin.Context) {

	path := c.Request.URL.Path
	method := c.Request.Method
	rateLimiterKey := method + ":" + path

	rateLimiter, ok := l.LimitRegistery[rateLimiterKey]

	if !ok {
		l.Logger.Error("Endpoint not registered", zap.String("endpoint", rateLimiterKey))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint not registered"})
		return
	}

	l.Logger.Info("Endpoint found", zap.String("endpoint", rateLimiterKey))

	strategyConfig := rateLimiter.Config
	allowed := rateLimiter.Limiter.Allow(c, strategyConfig)

	if !allowed {
		l.Logger.Error("Rate limit exceeded", zap.String("endpoint", rateLimiterKey))
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	targetUrl := config.Services[rateLimiter.TargetService]
	if targetUrl == "" {
		l.Logger.Error("Target service not configured", zap.String("service", rateLimiter.TargetService))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target service not configured"})
		return
	}

	proxy.ForwardRequest(c, targetUrl)
}
