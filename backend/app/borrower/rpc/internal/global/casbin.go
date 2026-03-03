package global

import (
	"os"

	gormadapter "github.com/casbin/gorm-adapter/v3"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/rbac"
	"zeebee-backend/domains/zeebee/constants"
)

var BorrowerEnforcer *rbac.Enforcer

func initCasbin(cnf config.CasbinConf) (err error) {
	policyContent, err := os.ReadFile(cnf.ModelPath)
	if err != nil {
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
	})
	if err != nil {
		return
	}
	err = BorrowerEnforcer.LoadFilteredPolicy(gormadapter.Filter{
		Ptype: []string{"p"},
		V0:    nil,
		V1:    []string{string(constants.UserTypeBorrower)},
		V2:    nil,
		V3:    nil,
		V4:    nil,
		V5:    nil,
	})

	return
}
