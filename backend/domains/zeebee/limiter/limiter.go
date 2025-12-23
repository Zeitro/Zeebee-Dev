package limiter

import (
	"context"
	"sync"

	"github.com/go-redis/redis_rate/v10"

	cpmmonlimiter "zeebee-backend/common/limiter"
)

type LimiterManager struct {
	limiter       *cpmmonlimiter.RedisLimiter
	strategyRates map[LimiterStrategy]redis_rate.Limit
	mux           sync.RWMutex
}

func NewLimiterManager(limiter *cpmmonlimiter.RedisLimiter, strategyRates map[LimiterStrategy]redis_rate.Limit) *LimiterManager {
	if strategyRates == nil {
		strategyRates = make(map[LimiterStrategy]redis_rate.Limit)
	}
	return &LimiterManager{
		limiter:       limiter,
		strategyRates: strategyRates,
		mux:           sync.RWMutex{},
	}
}

// AddStrategy
func (l *LimiterManager) AddStrategy(strategyKey LimiterStrategy, rate redis_rate.Limit) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.strategyRates[strategyKey] = rate
}

// RemoveStrategy
func (l *LimiterManager) RemoveStrategy(strategy LimiterStrategy) {
	l.mux.Lock()
	defer l.mux.Unlock()
	delete(l.strategyRates, strategy)
}

// GetStrategy
func (l *LimiterManager) GetStrategy(strategy LimiterStrategy) (rate redis_rate.Limit, isExisted bool) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	rate, isExisted = l.strategyRates[strategy]
	return rate, isExisted
}

func (l *LimiterManager) Limit(ctx context.Context, strategy LimiterStrategy, key string) (bool, error) {
	rate, isExisted := l.GetStrategy(strategy)
	if !isExisted {
		return false, nil
	}
	return l.limiter.Limit(ctx, key, rate)
}
