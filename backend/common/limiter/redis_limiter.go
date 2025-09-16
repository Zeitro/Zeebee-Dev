package limiter

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"

	"zeebee-backend/app/borrower/config"
)

type RedisLimiter struct {
	client redis_rate.Limiter
}

func GetLimiterRedisClient(conf config.RedisConf) (*redis.Client, error) {
	var tlsConfig *tls.Config
	if conf.Tls {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:      conf.Host,
		Password:  conf.Pass,
		DB:        conf.DB,
		TLSConfig: tlsConfig,
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func WithRedisMinuteLimit(number int) redis_rate.Limit {
	return redis_rate.PerMinute(number)
}

func WithRedisLimit(rate int, period time.Duration) redis_rate.Limit {
	return redis_rate.Limit{
		Rate:   rate,
		Burst:  rate,
		Period: period,
	}
}

func NewRedisLimiterWithRedis(rdb *redis.Client) *RedisLimiter {
	return &RedisLimiter{
		client: *redis_rate.NewLimiter(rdb),
	}
}

func (l *RedisLimiter) Limit(ctx context.Context, key string, limit redis_rate.Limit) (bool, error) {
	res, err := l.client.Allow(ctx, key, limit)
	if err != nil {
		return false, err
	}

	return res.Allowed > 0, nil
}
