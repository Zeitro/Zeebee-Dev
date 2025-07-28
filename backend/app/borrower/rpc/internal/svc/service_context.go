package svc

import (
	"zeebee-backend/app/borrower/rpc/internal/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.Config

	// repo
	UserRepo      *repo.UserRepo
	BorrowerRepo  *repo.BorrowerRepo
	LenderRepo    *repo.LenderRepo
	LoanRepo      *repo.LoanRepo
	PropertyRepo  *repo.PropertyRepo
	ApplicantRepo *repo.ApplicantRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// repo
		UserRepo:      repo.NewUserRepo(global.GetGormDB()),
		BorrowerRepo:  repo.NewBorrowerRepo(global.GetGormDB()),
		LenderRepo:    repo.NewLenderRepo(global.GetGormDB()),
		LoanRepo:      repo.NewLoanRepo(global.GetGormDB()),
		PropertyRepo:  repo.NewPropertyRepo(global.GetGormDB()),
		ApplicantRepo: repo.NewApplicantRepo(global.GetGormDB()),
	}
}
