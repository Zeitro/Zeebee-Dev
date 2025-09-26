package svc

import (
	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.RpcConf

	// repo
	UserRepo           *repo.UserRepo
	LenderCompanyRepo  *repo.LenderCompanyRepo
	LenderEmployeeRepo *repo.LenderEmployeeRepo
	UsGeographyRepo    *repo.UsGeographyRepo
	LoanRepo           *repo.LoanRepo
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		UserRepo:           repo.NewUserRepo(global.GetGormDB()),
		LenderCompanyRepo:  repo.NewLenderCompanyRepo(global.GetGormDB()),
		LenderEmployeeRepo: repo.NewLenderEmployeeRepo(global.GetGormDB()),
		UsGeographyRepo:    repo.NewUsGeographyRepo(global.GetGormDB()),
		LoanRepo:           repo.NewLoanRepo(global.GetGormDB()),
	}
}
