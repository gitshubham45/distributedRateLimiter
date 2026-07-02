package database

import (
	"fmt"
	"rate-limiter/internal/config"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func Connect(cfg config.DBConfig) {
	dbOnce.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			panic("Failed to connect to database")
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic("Failed to retrieve SQL database: " + err.Error())
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTimeSeconds) * time.Second)

		if err = sqlDB.Ping(); err != nil {
			panic("Failed to ping database: " + err.Error())
		}

		fmt.Println("Database connection established")
	})
}

func NewDB(cfg config.DBConfig) *gorm.DB {
	Connect(cfg)
	return db
}
