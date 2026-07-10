package svc

import (
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Logger logx.Logger
	Config config.RpcConf
	*repo.AllRepos
}

func NewServiceContext(c config.RpcConf, logger logx.Logger) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Logger:   logger,
		AllRepos: repo.NewAllRepos(global.GetGormDB()),
	}
}
