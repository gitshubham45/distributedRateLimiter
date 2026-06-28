package strategy

import (
	"rate-limiter/config"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type requestWindow struct {
	lastResetAt int64
	count       int
}

type PerUnitTimeLimit struct {
	mu            sync.Mutex
	requestWindow map[string]requestWindow
}

func NewPerUnitTimeLimit() *PerUnitTimeLimit {
	return &PerUnitTimeLimit{
		requestWindow: make(map[string]requestWindow),
	}
}

func (l *PerUnitTimeLimit) Allow(c *gin.Context, config config.StrategyConfig) bool {
	clientID := c.Query("client_id")
	key := c.Request.Method + ":" + c.Request.URL.Path + ":" + clientID

	now := time.Now().UnixNano()
	windowSize := config.Interval * int64(1e9)

	l.mu.Lock()
	defer l.mu.Unlock()

	w, found := l.requestWindow[key]
	if !found {
		l.requestWindow[key] = requestWindow{
			lastResetAt: now,
			count:       1,
		}
		return true
	}

	if now-w.lastResetAt >= windowSize {
		w.lastResetAt = now
		w.count = 1
		l.requestWindow[key] = w
		return true
	}

	if w.count < int(config.Limit) {
		w.count++
		l.requestWindow[key] = w
		return true
	}

	return false
}
