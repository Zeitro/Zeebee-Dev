package svc

import (
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Config config.RpcConf
	*repo.AllRepos

	// // repo
	// UserRepo *repo.UserRepo
	// // BorrowerRepo          *repo.BorrowerRepo
	// BorrowerEntityRepo    *repo.BorrowerEntityRepo
	// BorrowerApplicantRepo *repo.BorrowerApplicantRepo
	// // LenderCompanyRepo     *repo.LenderCompanyRepo
	// OrganizationRepo     *repo.OrganizationRepo
	// LoanRepo             *repo.LoanRepo
	// LoanPropertyRepo     *repo.LoanPropertyRepo
	// LoanBorrowerRepo     *repo.LoanBorrowerRepo
	// LoanDemographicRepo  *repo.LoanDemographicRepo
	// LoanExperienceRepo   *repo.LoanExperienceRepo
	// LoanThirdServiceRepo *repo.LoanThirdServiceRepo
	// LoanOperationRepo    *repo.LoanOperationRepo
	// LoanDocumentRepo     *repo.LoanDocumentRepo
	// BorrowerPropertyRepo *repo.BorrowerPropertyRepo
	// LoanApplicantRepo    *repo.LoanApplicantRepo
	// UsGeographyRepo      *repo.UsGeographyRepo
	// FileRepo             *repo.FileRepo
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		AllRepos: repo.NewAllRepos(global.GetGormDB()),

		// // repo
		// UserRepo: repo.NewUserRepo(global.GetGormDB()),
		// // BorrowerRepo:          repo.NewBorrowerRepo(global.GetGormDB()),
		// BorrowerEntityRepo:    repo.NewBorrowerEntityRepo(global.GetGormDB()),
		// BorrowerApplicantRepo: repo.NewBorrowerApplicantRepo(global.GetGormDB()),
		// // LenderCompanyRepo:     repo.NewLenderCompanyRepo(global.GetGormDB()),
		// OrganizationRepo:     repo.NewOrganizationRepo(global.GetGormDB()),
		// LoanRepo:             repo.NewLoanRepo(global.GetGormDB()),
		// LoanDemographicRepo:  repo.NewLoanDemographicRepo(global.GetGormDB()),
		// LoanExperienceRepo:   repo.NewLoanExperienceRepo(global.GetGormDB()),
		// LoanThirdServiceRepo: repo.NewLoanThirdServiceRepo(global.GetGormDB()),
		// LoanOperationRepo:    repo.NewLoanOperationRepo(global.GetGormDB()),
		// LoanBorrowerRepo:     repo.NewLoanBorrowerRepo(global.GetGormDB()),
		// LoanDocumentRepo:     repo.NewLoanDocumentRepo(global.GetGormDB()),
		// BorrowerPropertyRepo: repo.NewBorrowerPropertyRepo(global.GetGormDB()),
		// LoanPropertyRepo:     repo.NewLoanPropertyRepo(global.GetGormDB()),
		// LoanApplicantRepo:    repo.NewLoanApplicantRepo(global.GetGormDB()),
		// UsGeographyRepo:      repo.NewUsGeographyRepo(global.GetGormDB()),
		// FileRepo:             repo.NewFileRepo(global.GetGormDB()),
	}
}
