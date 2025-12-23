package svc

import (
	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Config config.RpcConf
	*repo.AllRepos
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		AllRepos: repo.NewAllRepos(global.GetGormDB()),
	}
}
