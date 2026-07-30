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

func initCasbin(cnf config.CasbinConf) (err error) {
	policyContent, err := os.ReadFile(cnf.ModelPath)
	if err != nil {
		err = errorx.Wrap(err, "read casbin model conf file error")
		return
	}

	AdminEnforcer, err = rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          GetGormDB(),
		ModelPolicy: string(policyContent),
		Watcher: rbac.WatcherConfig{
			RedisHost: cnf.Watcher.Host,
			Pwd:       cnf.Watcher.Pwd,
			Channel:   cnf.Watcher.AdminChannel,
			Tls:       cnf.Watcher.Tls,
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
		err = errorx.Wrap(err, "create admin rbac enforcer error")
		return
	}
	return
}
