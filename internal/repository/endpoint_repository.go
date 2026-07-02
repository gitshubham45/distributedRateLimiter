package repository

import (
	"context"
	"rate-limiter/internal/domain"
)

type EndpointRepository interface {
	Save(ctx context.Context, endpoint domain.EndpointConfig) error
	FindByKey(ctx context.Context, key string) (domain.EndpointConfig, error)
	FindAll(ctx context.Context) ([]domain.EndpointConfig, error)
}
