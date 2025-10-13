package svc

import (
	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.RpcConf

	// repo
	UserRepo             *repo.UserRepo
	LenderCompanyRepo    *repo.LenderCompanyRepo
	LenderEmployeeRepo   *repo.LenderEmployeeRepo
	UsGeographyRepo      *repo.UsGeographyRepo
	LoanRepo             *repo.LoanRepo
	BorrowerRepo         *repo.BorrowerRepo
	LoanPropertyRepo     *repo.LoanPropertyRepo
	LoanBorrowerRepo     *repo.LoanBorrowerRepo
	LoanDemographicRepo  *repo.LoanDemographicRepo
	LoanExperienceRepo   *repo.LoanExperienceRepo
	LoanThirdServiceRepo *repo.LoanThirdServiceRepo
	PropertyBaseRepo     *repo.PropertyBaseRepo
	ApplicantRepo        *repo.ApplicantRepo
	FileRepo             *repo.FileRepo
	LoanDocumentRepo     *repo.LoanDocumentRepo
	LoanOperationRepo    *repo.LoanOperationRepo
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		UserRepo:             repo.NewUserRepo(global.GetGormDB()),
		LenderCompanyRepo:    repo.NewLenderCompanyRepo(global.GetGormDB()),
		LenderEmployeeRepo:   repo.NewLenderEmployeeRepo(global.GetGormDB()),
		UsGeographyRepo:      repo.NewUsGeographyRepo(global.GetGormDB()),
		LoanRepo:             repo.NewLoanRepo(global.GetGormDB()),
		BorrowerRepo:         repo.NewBorrowerRepo(global.GetGormDB()),
		LoanPropertyRepo:     repo.NewLoanPropertyRepo(global.GetGormDB()),
		LoanBorrowerRepo:     repo.NewLoanBorrowerRepo(global.GetGormDB()),
		LoanDemographicRepo:  repo.NewLoanDemographicRepo(global.GetGormDB()),
		LoanExperienceRepo:   repo.NewLoanExperienceRepo(global.GetGormDB()),
		LoanThirdServiceRepo: repo.NewLoanThirdServiceRepo(global.GetGormDB()),
		PropertyBaseRepo:     repo.NewPropertyBaseRepo(global.GetGormDB()),
		ApplicantRepo:        repo.NewApplicantRepo(global.GetGormDB()),
		FileRepo:             repo.NewFileRepo(global.GetGormDB()),
		LoanDocumentRepo:     repo.NewLoanDocumentRepo(global.GetGormDB()),
		LoanOperationRepo:    repo.NewLoanOperationRepo(global.GetGormDB()),
	}
}
