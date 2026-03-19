package global

import (
	"os"

	gormadapter "github.com/casbin/gorm-adapter/v3"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/rbac"
	"zeebee-backend/domains/zeebee/constants"
)

var AdminEnforcer *rbac.Enforcer
var LenderEnforcer *rbac.Enforcer
var BorrowerEnforcer *rbac.Enforcer

func initCasbin(cnf config.CasbinConf) (err error) {
	policyContent, err := os.ReadFile(cnf.ModelPath)
	if err != nil {
		err = errorx.Wrap(err, "Read casbin model conf file error")
		return
	}

	AdminEnforcer, err = rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          GetGormDB(),
		ModelPolicy: string(policyContent),
		Watcher: rbac.WatcherConfig{
			RedisHost: cnf.Watcher.Host,
			Pwd:       cnf.Watcher.Pwd,
			Channel:   cnf.Watcher.AdminChannel,
		},
		DefaultLoadFilteredPolicy: &gormadapter.Filter{
			Ptype: []string{"p"},
			V0:    nil,
			V1:    []string{string(constants.UserTypeAdmin)},
			V2:    nil,
			V3:    nil,
			V4:    nil,
			V5:    nil,
		},
	})
	if err != nil {
		err = errorx.Wrap(err, "Create admin rbac enforcer error")
		return
	}

	LenderEnforcer, err = rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          GetGormDB(),
		ModelPolicy: string(policyContent),
		Watcher: rbac.WatcherConfig{
			RedisHost: cnf.Watcher.Host,
			Pwd:       cnf.Watcher.Pwd,
			Channel:   cnf.Watcher.LenderChannel,
		},
		DefaultLoadFilteredPolicy: &gormadapter.Filter{
			Ptype: []string{"p"},
			V0:    nil,
			V1:    []string{string(constants.UserTypeLender)},
			V2:    nil,
			V3:    nil,
			V4:    nil,
			V5:    nil,
		},
	})
	if err != nil {
		err = errorx.Wrap(err, "Create lender rbac enforcer error")
		return
	}

	BorrowerEnforcer, err = rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          GetGormDB(),
		ModelPolicy: string(policyContent),
		Watcher: rbac.WatcherConfig{
			RedisHost: cnf.Watcher.Host,
			Pwd:       cnf.Watcher.Pwd,
			Channel:   cnf.Watcher.BorrowerChannel,
		},
		DefaultLoadFilteredPolicy: &gormadapter.Filter{
			Ptype: []string{"p"},
			V0:    nil,
			V1:    []string{string(constants.UserTypeBorrower)},
			V2:    nil,
			V3:    nil,
			V4:    nil,
			V5:    nil,
		},
	})
	if err != nil {
		err = errorx.Wrap(err, "Create borrower rbac enforcer error")
		return
	}

	return
}
