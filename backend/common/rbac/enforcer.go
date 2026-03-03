package rbac

import (
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
)

type Enforcer struct {
	Enforcer *casbin.Enforcer
}

type EnforcerConfig struct {
	DB          *gorm.DB
	ModelPolicy string
	Watcher     WatcherConfig
}

type WatcherConfig struct {
	RedisHost string
	Pwd       string
	Channel   string
}

func NewEnforcer(config EnforcerConfig) (*Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(config.DB)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	casbinModel, err := model.NewModelFromString(config.ModelPolicy)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	enforcer, err := casbin.NewEnforcer(casbinModel, adapter)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	enforcer.EnableAutoSave(true)

	// redis watch
	redisHost := strings.TrimSpace(config.Watcher.RedisHost)
	if len(redisHost) > 0 {
		var watcher persist.Watcher
		watcher, err = rediswatcher.NewWatcher(redisHost, rediswatcher.WatcherOptions{
			Channel:    config.Watcher.Channel,
			IgnoreSelf: true,
			Options: redis.Options{
				Network:  "tcp",
				Password: config.Watcher.Pwd,
			},
		})
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		err = watcher.SetUpdateCallback(rediswatcher.DefaultUpdateCallback(enforcer))
		if err != nil {
			return nil, errorx.WithStack(err)
		}

		err = enforcer.SetWatcher(watcher)
		if err != nil {
			return nil, errorx.WithStack(err)
		}

	}

	return &Enforcer{
		Enforcer: enforcer,
	}, nil
}

func (e *Enforcer) LoadFilteredPolicy(filter gormadapter.Filter) error {
	err := e.Enforcer.LoadFilteredPolicy(filter)
	return errorx.WithStack(err)
}

func (e *Enforcer) Enforce(rvals ...interface{}) (bool, error) {
	result, err := e.Enforcer.Enforce(rvals...)
	return result, errorx.WithStack(err)
}

// AddPolicies
func (e *Enforcer) AddPolicies(p [][]string) error {
	_, err := e.Enforcer.AddPolicies(p)
	return errorx.WithStack(err)
}

// RemovePolicies
func (e *Enforcer) RemovePolicies(p [][]string) error {
	_, err := e.Enforcer.RemovePolicies(p)
	return errorx.WithStack(err)
}

func (e *Enforcer) SavePolicy() error {
	err := e.Enforcer.SavePolicy()
	return errorx.WithStack(err)
}
