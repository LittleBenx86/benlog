package dependencies

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Dependencies struct {
	DBClient    *gorm.DB
	RedisClient *redis.Client
	Logger      *zap.Logger
}

type Option func(dependencies *Dependencies)

func WithDBClient(client *gorm.DB) Option {
	return func(d *Dependencies) {
		d.DBClient = client
	}
}

func WithRedisClient(client *redis.Client) Option {
	return func(d *Dependencies) {
		d.RedisClient = client
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(d *Dependencies) {
		d.Logger = logger
	}
}
