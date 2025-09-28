package svc

import (
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.RpcConf

	// repo
	UserRepo             *repo.UserRepo
	BorrowerRepo         *repo.BorrowerRepo
	LenderRepo           *repo.LenderRepo
	LoanRepo             *repo.LoanRepo
	LoanPropertyRepo     *repo.LoanPropertyRepo
	LoanBorrowerRepo     *repo.LoanBorrowerRepo
	LoanDemographicRepo  *repo.LoanDemographicRepo
	LoanExperienceRepo   *repo.LoanExperienceRepo
	LoanThirdServiceRepo *repo.LoanThirdServiceRepo
	LoanOperationRepo    *repo.LoanOperationRepo
	PropertyBaseRepo     *repo.PropertyBaseRepo
	ApplicantRepo        *repo.ApplicantRepo
	UsGeographyRepo      *repo.UsGeographyRepo
	FileRepo             *repo.FileRepo
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// repo
		UserRepo:             repo.NewUserRepo(global.GetGormDB()),
		BorrowerRepo:         repo.NewBorrowerRepo(global.GetGormDB()),
		LenderRepo:           repo.NewLenderRepo(global.GetGormDB()),
		LoanRepo:             repo.NewLoanRepo(global.GetGormDB()),
		LoanDemographicRepo:  repo.NewLoanDemographicRepo(global.GetGormDB()),
		LoanExperienceRepo:   repo.NewLoanExperienceRepo(global.GetGormDB()),
		LoanThirdServiceRepo: repo.NewLoanThirdServiceRepo(global.GetGormDB()),
		LoanOperationRepo:    repo.NewLoanOperationRepo(global.GetGormDB()),
		LoanBorrowerRepo:     repo.NewLoanBorrowerRepo(global.GetGormDB()),
		PropertyBaseRepo:     repo.NewPropertyBaseRepo(global.GetGormDB()),
		LoanPropertyRepo:     repo.NewLoanPropertyRepo(global.GetGormDB()),
		ApplicantRepo:        repo.NewApplicantRepo(global.GetGormDB()),
		UsGeographyRepo:      repo.NewUsGeographyRepo(global.GetGormDB()),
		FileRepo:             repo.NewFileRepo(global.GetGormDB()),
	}
}
