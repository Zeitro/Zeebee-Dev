package svc

import (
	"zeebee-backend/app/lender/config"
)

type ServiceContext struct {
	Config config.ApiConf
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
