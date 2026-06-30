package limiter

import (
	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter/strategy"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	strategyRegistry = make(map[string]RequestLimiter)
	strategyMutex    sync.RWMutex
)

type RequestLimiter interface {
	Allow(c *gin.Context, config config.StrategyConfig) bool
}

func GetRequestLimiter(name string) RequestLimiter {
	strategyMutex.RLock()
	requestLimiter, found := strategyRegistry[name]
	strategyMutex.RUnlock()

	if found {
		return requestLimiter
	}

	strategyMutex.Lock()
	defer strategyMutex.Unlock()

	if requestLimiter, found := strategyRegistry[name]; found {
		return requestLimiter
	}

	switch name {
	case config.PER_UNIT_TIME:
		requestLimiter = strategy.NewPerUnitTimeLimit()
	case config.LEAKY_BUCKET:
		requestLimiter = &strategy.LeakyBucket{}
	default:
		requestLimiter = strategy.NewPerUnitTimeLimit()
	}

	strategyRegistry[name] = requestLimiter
	return requestLimiter
}
