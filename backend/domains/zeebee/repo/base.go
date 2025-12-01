package repo

import "gorm.io/gorm"

type AllRepos struct {
	BorrowerApplicantRepo   *BorrowerApplicantRepo
	BorrowerEntityRepo      *BorrowerEntityRepo
	BorrowerPropertyRepo    *BorrowerPropertyRepo
	FileRepo                *FileRepo
	LoanApplicantRepo       *LoanApplicantRepo
	LoanBorrowerRepo        *LoanBorrowerRepo
	LoanDemographicRepo     *LoanDemographicRepo
	LoanDocumentRepo        *LoanDocumentRepo
	LoanExperienceRepo      *LoanExperienceRepo
	LoanOperationRepo       *LoanOperationRepo
	LoanPropertyRepo        *LoanPropertyRepo
	LoanRepo                *LoanRepo
	LoanThirdServiceRepo    *LoanThirdServiceRepo
	OrganizationRepo        *OrganizationRepo
	OrganizationAndUserRepo *OrganizationAndUserRepo
	UsGeographyRepo         *UsGeographyRepo
	UserRepo                *UserRepo
	ApiRepo                 *ApiRepo
	ResourceRepo            *ResourceRepo
	RoleRepo                *RoleRepo
}

func NewAllRepos(db *gorm.DB) *AllRepos {
	return &AllRepos{
		BorrowerApplicantRepo:   NewBorrowerApplicantRepo(db),
		BorrowerEntityRepo:      NewBorrowerEntityRepo(db),
		BorrowerPropertyRepo:    NewBorrowerPropertyRepo(db),
		FileRepo:                NewFileRepo(db),
		LoanApplicantRepo:       NewLoanApplicantRepo(db),
		LoanBorrowerRepo:        NewLoanBorrowerRepo(db),
		LoanDemographicRepo:     NewLoanDemographicRepo(db),
		LoanDocumentRepo:        NewLoanDocumentRepo(db),
		LoanExperienceRepo:      NewLoanExperienceRepo(db),
		LoanOperationRepo:       NewLoanOperationRepo(db),
		LoanPropertyRepo:        NewLoanPropertyRepo(db),
		LoanRepo:                NewLoanRepo(db),
		LoanThirdServiceRepo:    NewLoanThirdServiceRepo(db),
		OrganizationRepo:        NewOrganizationRepo(db),
		OrganizationAndUserRepo: NewOrganizationAndUserRepo(db),
		UsGeographyRepo:         NewUsGeographyRepo(db),
		UserRepo:                NewUserRepo(db),
		ApiRepo:                 NewApiRepo(db),
		ResourceRepo:            NewResourceRepo(db),
		RoleRepo:                NewRoleRepo(db),
	}
}
