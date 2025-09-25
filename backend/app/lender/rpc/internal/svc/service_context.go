package svc

import "zeebee-backend/app/lender/config"

type ServiceContext struct {
	Config config.RpcConf
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
