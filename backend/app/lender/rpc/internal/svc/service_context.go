package svc

import (
	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.RpcConf

	// repo
	UsGeographyRepo *repo.UsGeographyRepo
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		UsGeographyRepo: repo.NewUsGeographyRepo(global.GetGormDB()),
	}
}
