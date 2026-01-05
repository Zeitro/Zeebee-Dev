package svc

import (
	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Config config.RpcConf
	*repo.AllRepos
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		AllRepos: repo.NewAllRepos(global.GetGormDB()),
	}
}
