package svc

import (
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.RpcConf

	// repo
	UserRepo        *repo.UserRepo
	BorrowerRepo    *repo.BorrowerRepo
	LenderRepo      *repo.LenderRepo
	LoanRepo        *repo.LoanRepo
	PropertyRepo    *repo.PropertyRepo
	ApplicantRepo   *repo.ApplicantRepo
	UsGeographyRepo *repo.UsGeographyRepo
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// repo
		UserRepo:        repo.NewUserRepo(global.GetGormDB()),
		BorrowerRepo:    repo.NewBorrowerRepo(global.GetGormDB()),
		LenderRepo:      repo.NewLenderRepo(global.GetGormDB()),
		LoanRepo:        repo.NewLoanRepo(global.GetGormDB()),
		PropertyRepo:    repo.NewPropertyRepo(global.GetGormDB()),
		ApplicantRepo:   repo.NewApplicantRepo(global.GetGormDB()),
		UsGeographyRepo: repo.NewUsGeographyRepo(global.GetGormDB()),
	}
}
