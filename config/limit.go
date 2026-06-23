package config

import (
	"rate-limiter/limiter"
)

func InitStrategyMap() map[string]limiter.RequestLimiter {
	strategyMap := make(map[string]limiter.RequestLimiter)
	strategyMap["per_hour"] = &limiter.PerHourLimit{}
	strategyMap["per_min"]
	return strategyMap
}
