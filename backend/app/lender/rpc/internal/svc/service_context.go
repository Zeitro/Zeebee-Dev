package svc

import (
	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Config config.RpcConf
	*repo.AllRepos

	// repo
	// UserRepo              *repo.UserRepo
	// OrganizationRepo      *repo.OrganizationRepo
	// UsGeographyRepo       *repo.UsGeographyRepo
	// LoanRepo              *repo.LoanRepo
	// BorrowerRepo          *repo.BorrowerRepo
	// LoanPropertyRepo      *repo.LoanPropertyRepo
	// LoanBorrowerRepo      *repo.LoanBorrowerRepo
	// LoanDemographicRepo   *repo.LoanDemographicRepo
	// LoanExperienceRepo    *repo.LoanExperienceRepo
	// LoanThirdServiceRepo  *repo.LoanThirdServiceRepo
	// BorrowerPropertyRepo  *repo.BorrowerPropertyRepo
	// BorrowerEntityRepo    *repo.BorrowerEntityRepo
	// BorrowerApplicantRepo *repo.BorrowerApplicantRepo
	// LoanApplicantRepo     *repo.LoanApplicantRepo
	// FileRepo              *repo.FileRepo
	// LoanDocumentRepo      *repo.LoanDocumentRepo
	// LoanOperationRepo     *repo.LoanOperationRepo
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		AllRepos: repo.NewAllRepos(global.GetGormDB()),
	}
}
