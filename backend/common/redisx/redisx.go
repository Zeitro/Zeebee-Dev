package redisx

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"zeebee-backend/common/errorx"
)

type Redisx struct {
	rds    *redis.Client
	prefix string
}

type Opt struct {
	Host string
	Type string
	Pass string
	Tls  bool
	DB   int

	Prefix string
}

func GetRedis(opt Opt) (rds *Redisx, err error) {
	var rdsCli *redis.Client
	redisOpt := &redis.Options{
		Addr:     opt.Host,
		Password: opt.Pass,
		DB:       opt.DB,
	}
	if opt.Tls {
		redisOpt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	rdsCli = redis.NewClient(redisOpt)

	ctx := context.Background()
	if err = rdsCli.Ping(ctx).Err(); err != nil {
		return nil, errorx.Wrap(err, "Redis ping failed")
	}

	return &Redisx{
		rds:    rdsCli,
		prefix: opt.Prefix,
	}, nil
}

func (r *Redisx) GetClient() *redis.Client {
	return r.rds
}

func (r *Redisx) Get(ctx context.Context, key string) (val string, err error) {
	return r.rds.Get(ctx, r.buildKey(key)).Result()
}

func (r *Redisx) GetInt(ctx context.Context, key string) (val int, err error) {
	strVal, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if len(strVal) == 0 {
		return 0, nil
	}
	return strconv.Atoi(strVal)
}

// Set
func (r *Redisx) Set(ctx context.Context, key string, val string) (err error) {
	return r.rds.Set(ctx, r.buildKey(key), val, 0).Err()
}

// SetExpire
func (r *Redisx) SetExpire(ctx context.Context, key string, val string, seconds int) (err error) {
	return r.rds.Set(ctx, r.buildKey(key), val, time.Duration(seconds)*time.Second).Err()
}

// Expire
func (r *Redisx) Expire(ctx context.Context, key string, seconds int) (err error) {
	return r.rds.Expire(ctx, r.buildKey(key), time.Duration(seconds)*time.Second).Err()
}

// Increase
func (r *Redisx) Increase(ctx context.Context, key string) (newVal int64, err error) {
	return r.rds.Incr(ctx, r.buildKey(key)).Result()
}

// IncreaseExpire
func (r *Redisx) IncreaseExpire(ctx context.Context, key string, seconds int) (newVal int64, err error) {
	newVal, err = r.rds.Incr(ctx, r.buildKey(key)).Result()
	if err != nil {
		return 0, err
	}
	if newVal == 1 {
		err = r.Expire(ctx, key, seconds)
	}
	return newVal, err
}

// IncreaseBy
func (r *Redisx) IncreaseBy(ctx context.Context, key string, val int64) (newVal int64, err error) {
	return r.rds.IncrBy(ctx, r.buildKey(key), val).Result()
}

// Del
func (r *Redisx) Del(ctx context.Context, key string) (err error) {
	return r.rds.Del(ctx, r.buildKey(key)).Err()
}

// IsExisted
func (r *Redisx) IsExisted(ctx context.Context, key string) (bool, error) {
	count, err := r.rds.Exists(ctx, r.buildKey(key)).Result()
	return count > 0, err
}

func (r *Redisx) buildKey(key string) string {
	if len(r.prefix) == 0 || (len(r.prefix) > 0 && strings.HasPrefix(key, r.prefix+":")) {
		return key
	}
	return fmt.Sprintf("%s:%s", r.prefix, key)
}
