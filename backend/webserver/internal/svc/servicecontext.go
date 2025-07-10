package svc

import (
	"zeebee-backend/webserver/internal/config"
	"zeebee-backend/webserver/internal/model"
)

type ServiceContext struct {
	Config    config.Config
	UserModel *model.UserModel
	LoanModel *model.LoanModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(c.DataSource),
		LoanModel: model.NewLoanModel(c.DataSource),
	}
}
