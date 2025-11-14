package redisx

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Redisx struct {
	rds    *redis.Redis
	prefix string
}

type Opt struct {
	Host string
	Type string
	Pass string
	Tls  bool

	Prefix string
}

func GetRedis(opt Opt) (rds *Redisx, err error) {
	var rdsCli *redis.Redis
	rdsCli, err = redis.NewRedis(redis.RedisConf{
		Host: opt.Host,
		Type: opt.Type,
		Pass: opt.Pass,
		Tls:  opt.Tls,
	})
	if err != nil {
		return nil, errors.Wrap(err, "New Redis Error")
	}

	if !rdsCli.Ping() {
		return nil, errors.New("Redis Ping Failed")
	}
	return &Redisx{
		rds:    rdsCli,
		prefix: opt.Prefix,
	}, nil
}

func (r *Redisx) GetClient() *redis.Redis {
	return r.rds
}

func (r *Redisx) Get(ctx context.Context, key string) (val string, err error) {
	return r.rds.GetCtx(ctx, r.buildKey(key))
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
	return r.rds.SetCtx(ctx, r.buildKey(key), val)
}

// SetExpire
func (r *Redisx) SetExpire(ctx context.Context, key string, val string, seconds int) (err error) {
	return r.rds.SetexCtx(ctx, r.buildKey(key), val, seconds)
}

// Expire
func (r *Redisx) Expire(ctx context.Context, key string, seconds int) (err error) {
	return r.rds.ExpireCtx(ctx, r.buildKey(key), seconds)
}

// Increase
func (r *Redisx) Increase(ctx context.Context, key string) (newVal int64, err error) {
	return r.rds.IncrCtx(ctx, r.buildKey(key))
}

// IncreaseExpire
func (r *Redisx) IncreaseExpire(ctx context.Context, key string, seconds int) (newVal int64, err error) {
	newVal, err = r.rds.IncrCtx(ctx, r.buildKey(key))
	if err != nil {
		return 0, err
	}
	if newVal == 1 {
		err = r.Expire(ctx, r.buildKey(key), seconds)
	}
	return newVal, err
}

// IncreaseBy
func (r *Redisx) IncreaseBy(ctx context.Context, key string, val int64) (newVal int64, err error) {
	return r.rds.IncrbyCtx(ctx, r.buildKey(key), val)
}

// Del
func (r *Redisx) Del(ctx context.Context, key string) (err error) {
	_, err = r.rds.DelCtx(ctx, r.buildKey(key))
	return
}

// IsExisted
func (r *Redisx) IsExisted(ctx context.Context, key string) (bool, error) {
	return r.rds.ExistsCtx(ctx, r.buildKey(key))
}

func (r *Redisx) buildKey(key string) string {
	if len(r.prefix) == 0 || (len(r.prefix) > 0 && strings.HasPrefix(key, r.prefix+":")) {
		return key
	}
	return fmt.Sprintf("%s:%s", r.prefix, key)
}
