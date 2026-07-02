package repository

import (
	"context"
	"rate-limiter/internal/domain"
	"time"

	"gorm.io/gorm"
)

type EndpointEntity struct {
	ID            uint   `gorm:"primaryKey"`
	Key           string `gorm:"uniqueIndex;not null"`
	Method        string `gorm:"not null"`
	Path          string `gorm:"not null"`
	TargetService string `gorm:"not null"`
	StrategyName  string `gorm:"not null"`
	Limit         int64  `gorm:"column:limit_value;not null"`
	Interval      int64  `gorm:"column:interval_seconds;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (EndpointEntity) TableName() string {
	return "endpoints"
}

type PostgresEndpointRepository struct {
	db *gorm.DB
}

func NewPostgresEndpointRepository(db *gorm.DB) *PostgresEndpointRepository {
	return &PostgresEndpointRepository{
		db: db,
	}
}

func (r *PostgresEndpointRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&EndpointEntity{})
}

func (r *PostgresEndpointRepository) Save(ctx context.Context, endpoint domain.EndpointConfig) error {
	entity := toEntity(endpoint)
	return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *PostgresEndpointRepository) FindByKey(ctx context.Context, key string) (domain.EndpointConfig, error) {
	var entity EndpointEntity
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&entity).Error; err != nil {
		return domain.EndpointConfig{}, err
	}

	return toDomain(entity), nil
}

func (r *PostgresEndpointRepository) FindAll(ctx context.Context) ([]domain.EndpointConfig, error) {
	var entities []EndpointEntity
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}

	endpoints := make([]domain.EndpointConfig, 0, len(entities))
	for _, entity := range entities {
		endpoints = append(endpoints, toDomain(entity))
	}

	return endpoints, nil
}

func toEntity(endpoint domain.EndpointConfig) EndpointEntity {
	return EndpointEntity{
		Key:           endpoint.Key,
		Method:        endpoint.Method,
		Path:          endpoint.Path,
		TargetService: endpoint.TargetService,
		StrategyName:  endpoint.StrategyName,
		Limit:         endpoint.Limit,
		Interval:      endpoint.Interval,
	}
}

func toDomain(entity EndpointEntity) domain.EndpointConfig {
	return domain.EndpointConfig{
		Key:           entity.Key,
		Method:        entity.Method,
		Path:          entity.Path,
		TargetService: entity.TargetService,
		StrategyName:  entity.StrategyName,
		Limit:         entity.Limit,
		Interval:      entity.Interval,
	}
}
