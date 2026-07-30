package svc

import (
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Config config.RpcConf
	Logger logx.Logger
	*repo.AllRepos
}

func NewServiceContext(c config.RpcConf, logger logx.Logger) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Logger:   logger,
		AllRepos: repo.NewAllRepos(global.GetGormDB()),
	}
}
