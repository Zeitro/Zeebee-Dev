package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client redis_rate.Limiter
}

func GetLimiterRedisClient(host string, pwd string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
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

func NewRedisLimiterWithHost(host string, pwd string, db int) (*RedisLimiter, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &RedisLimiter{
		client: *redis_rate.NewLimiter(rdb),
	}, nil
}

func (l *RedisLimiter) Limit(ctx context.Context, key string, limit redis_rate.Limit) (bool, error) {
	res, err := l.client.Allow(ctx, key, limit)
	if err != nil {
		return false, err
	}

	return res.Allowed > 0, nil
}
