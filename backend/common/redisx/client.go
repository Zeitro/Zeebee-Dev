package redisx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"zeebee-backend/common/errorx"
)

type Client struct {
	cli    *redis.Client
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

func GetRedis(opt Opt) (cli *Client, err error) {
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

	return &Client{
		cli:    rdsCli,
		prefix: opt.Prefix,
	}, nil
}

func (r *Client) GetClient() *redis.Client {
	return r.cli
}

func (r *Client) Get(ctx context.Context, key string) (val string, err error) {
	val, err = r.cli.Get(ctx, r.buildKey(key)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return val, err
}

func (r *Client) GetInt(ctx context.Context, key string) (val int, err error) {
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
func (r *Client) Set(ctx context.Context, key string, val string) (err error) {
	return r.cli.Set(ctx, r.buildKey(key), val, 0).Err()
}

// SetExpire
func (r *Client) SetExpire(ctx context.Context, key string, val string, seconds int) (err error) {
	return r.cli.Set(ctx, r.buildKey(key), val, time.Duration(seconds)*time.Second).Err()
}

// Expire
func (r *Client) Expire(ctx context.Context, key string, seconds int) (err error) {
	return r.cli.Expire(ctx, r.buildKey(key), time.Duration(seconds)*time.Second).Err()
}

// Increase
func (r *Client) Increase(ctx context.Context, key string) (newVal int64, err error) {
	return r.cli.Incr(ctx, r.buildKey(key)).Result()
}

// IncreaseExpire
func (r *Client) IncreaseExpire(ctx context.Context, key string, seconds int) (newVal int64, err error) {
	newVal, err = r.cli.Incr(ctx, r.buildKey(key)).Result()
	if err != nil {
		return 0, err
	}
	if newVal == 1 {
		err = r.Expire(ctx, key, seconds)
	}
	return newVal, err
}

// IncreaseBy
func (r *Client) IncreaseBy(ctx context.Context, key string, val int64) (newVal int64, err error) {
	return r.cli.IncrBy(ctx, r.buildKey(key), val).Result()
}

// Del
func (r *Client) Del(ctx context.Context, key string) (err error) {
	return r.cli.Del(ctx, r.buildKey(key)).Err()
}

// IsExisted
func (r *Client) IsExisted(ctx context.Context, key string) (bool, error) {
	count, err := r.cli.Exists(ctx, r.buildKey(key)).Result()
	return count > 0, err
}

func (r *Client) buildKey(key string) string {
	if len(r.prefix) == 0 || (len(r.prefix) > 0 && strings.HasPrefix(key, r.prefix+":")) {
		return key
	}
	return fmt.Sprintf("%s:%s", r.prefix, key)
}
