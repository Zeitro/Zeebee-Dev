package rbac

import (
	"crypto/tls"
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
	Enforcer                  *casbin.Enforcer
	defaultLoadFilteredPolicy *gormadapter.Filter
}

type EnforcerConfig struct {
	DB                        *gorm.DB
	ModelPolicy               string
	Watcher                   WatcherConfig
	DefaultLoadFilteredPolicy *gormadapter.Filter
}

type WatcherConfig struct {
	RedisHost string
	Pwd       string
	Tls       bool
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
		redisOption := redis.Options{
			Network:  "tcp",
			Password: config.Watcher.Pwd,
		}
		if config.Watcher.Tls {
			redisOption.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		watcher, err = rediswatcher.NewWatcher(redisHost, rediswatcher.WatcherOptions{
			Channel:    config.Watcher.Channel,
			IgnoreSelf: true,
			Options:    redisOption,
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

	e := &Enforcer{
		Enforcer:                  enforcer,
		defaultLoadFilteredPolicy: config.DefaultLoadFilteredPolicy,
	}
	err = e.LoadDefaultFilteredPolicy()
	return e, errorx.WithStack(err)
}

func (e *Enforcer) LoadDefaultFilteredPolicy() error {
	if e.defaultLoadFilteredPolicy == nil {
		return nil
	}
	err := e.Enforcer.LoadFilteredPolicy(*e.defaultLoadFilteredPolicy)
	return errorx.WithStack(err)
}

func (e *Enforcer) LoadFilteredPolicy(filter gormadapter.Filter) error {
	err := e.Enforcer.LoadFilteredPolicy(filter)
	return errorx.WithStack(err)
}

func (e *Enforcer) Enforce(rvals ...interface{}) (bool, error) {
	result, err := e.Enforcer.Enforce(rvals...)
	return result, errorx.WithStack(err)
}

// AddPolicies adds multiple policies
func (e *Enforcer) AddPolicies(p [][]string) error {
	_, err := e.Enforcer.AddPolicies(p)
	return errorx.WithStack(err)
}

// RemovePolicies removes multiple policies
func (e *Enforcer) RemovePolicies(p [][]string) error {
	_, err := e.Enforcer.RemovePolicies(p)
	return errorx.WithStack(err)
}

// SavePolicy saves all added policies to DB
func (e *Enforcer) SavePolicy() error {
	err := e.Enforcer.SavePolicy()
	return errorx.WithStack(err)
}

// AllPolicies returns all policies currently loaded in memory.
// Since LoadDefaultFilteredPolicy is called during initialization,
// this returns the filtered policies for a specific endpoint (admin/borrower/lender).
func (e *Enforcer) AllPolicies() (policies [][]string, err error) {
	policies, err = e.Enforcer.GetPolicy()
	return policies, err
}
