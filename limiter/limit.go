package limiter

import (
	"rate-limiter/config"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	strategyRegistry = make(map[string]RequestLimiter)
	strategyMutex    sync.RWMutex
)

type LimitConfig struct {
	limit    int64
	interval int64
}

type RequestLimiter interface {
	Allow(c *gin.Context, config LimitConfig) bool
}

type PerUnitTimeLimit struct{}

func (l *PerUnitTimeLimit) Allow(c *gin.Context, config LimitConfig) bool {
	return true
}

type LeakyBucket struct{}

func (l *LeakyBucket) Allow(c *gin.Context, config LimitConfig) bool {
	return true
}

func GetRequestLimiter(name string) RequestLimiter {
	switch name {
	case config.PER_UNIT_TIME:
		strategyMutex.RLock()
		defer strategyMutex.RUnlock()
		if _, found := strategyRegistry[name]; !found {
			strategyRegistry[name] = &PerUnitTimeLimit{}
		}
		return strategyRegistry[name]
	case config.LEAKY_BUCKET:
		strategyMutex.RLock()
		defer strategyMutex.RUnlock()

		if _, found := strategyRegistry[name]; !found {
			strategyRegistry[name] = &LeakyBucket{}
		}
		return strategyRegistry[name]
	default:
		strategyMutex.RLock()
		defer strategyMutex.RUnlock()

		if _, found := strategyRegistry[name]; !found {
			strategyRegistry[name] = &PerUnitTimeLimit{}
		}
		return strategyRegistry[name]
	}
}
