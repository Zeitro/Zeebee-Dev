package repo

import "gorm.io/gorm"

type AllRepos struct {
	BorrowerApplicantRepo           *BorrowerApplicantRepo
	BorrowerEntityRepo              *BorrowerEntityRepo
	BorrowerPropertyRepo            *BorrowerPropertyRepo
	FileRepo                        *FileRepo
	LoanApplicantRepo               *LoanApplicantRepo
	LoanBorrowerRepo                *LoanBorrowerRepo
	LoanDemographicRepo             *LoanDemographicRepo
	LoanDocumentRepo                *LoanDocumentRepo
	LoanExperienceRepo              *LoanExperienceRepo
	LoanOperationRepo               *LoanOperationRepo
	LoanPropertyRepo                *LoanPropertyRepo
	LoanRepo                        *LoanRepo
	LoanThirdServiceRepo            *LoanThirdServiceRepo
	OrganizationRepo                *OrganizationRepo
	OrganizationUserRepo            *OrganizationUserRepo
	UsGeographyRepo                 *UsGeographyRepo
	UserRepo                        *UserRepo
	ApiRepo                         *ApiRepo
	ResourceRepo                    *ResourceRepo
	RoleRepo                        *RoleRepo
	RoleAndResourceRepo             *RoleAndResourceRepo
	EmailAuthorizationRepo          *EmailAuthorizationRepo
	OrganizationEmailRepo           *OrganizationEmailRepo
	BaseTemplateRepo                *BaseTemplateRepo
	OrganizationTemplateRepo        *OrganizationTemplateRepo
	OrganizationSettingRepo         *OrganizationSettingRepo
	BorrowerPropertyTransactionRepo *BorrowerPropertyTransactionRepo
	OperationRepo                   *OperationRepo
	DocumentRepo                    *DocumentRepo
	EntityDocumentRepo              *EntityDocumentRepo
	DocumentTemplateRepo            *DocumentTemplateRepo
	DriveAuthorizationRepo          *DriveAuthorizationRepo
	DriveSyncJobRepo                *DriveSyncJobRepo
	PricingPlanRepo                 *PricingPlanRepo
	PricingPlanPolicyRepo           *PricingPlanPolicyRepo
	PricingPlanAndResourceRepo      *PricingPlanAndResourceRepo
	OrganizationSubscriptionRepo    *OrganizationSubscriptionRepo
	StripeEventRepo                 *StripeEventRepo
	PaymentTransactionRepo          *PaymentTransactionRepo
}

// NewAllRepos creates a new AllRepos instance with all repository implementations initialized.
// It takes a gorm.DB instance as a parameter and returns a pointer to AllRepos.
// Each repository is initialized with the provided database connection.
func NewAllRepos(db *gorm.DB) *AllRepos {
	return &AllRepos{
		BorrowerApplicantRepo:           NewBorrowerApplicantRepo(db),
		BorrowerEntityRepo:              NewBorrowerEntityRepo(db),
		BorrowerPropertyRepo:            NewBorrowerPropertyRepo(db),
		FileRepo:                        NewFileRepo(db),
		LoanApplicantRepo:               NewLoanApplicantRepo(db),
		LoanBorrowerRepo:                NewLoanBorrowerRepo(db),
		LoanDemographicRepo:             NewLoanDemographicRepo(db),
		LoanDocumentRepo:                NewLoanDocumentRepo(db),
		LoanExperienceRepo:              NewLoanExperienceRepo(db),
		LoanOperationRepo:               NewLoanOperationRepo(db),
		LoanPropertyRepo:                NewLoanPropertyRepo(db),
		LoanRepo:                        NewLoanRepo(db),
		LoanThirdServiceRepo:            NewLoanThirdServiceRepo(db),
		OrganizationRepo:                NewOrganizationRepo(db),
		OrganizationUserRepo:            NewOrganizationUserRepo(db),
		UsGeographyRepo:                 NewUsGeographyRepo(db),
		UserRepo:                        NewUserRepo(db),
		ApiRepo:                         NewApiRepo(db),
		ResourceRepo:                    NewResourceRepo(db),
		RoleRepo:                        NewRoleRepo(db),
		RoleAndResourceRepo:             NewRoleAndResourceRepo(db),
		EmailAuthorizationRepo:          NewEmailAuthorizationRepo(db),
		OrganizationEmailRepo:           NewOrganizationEmailRepo(db),
		BaseTemplateRepo:                NewBaseTemplateRepo(db),
		OrganizationTemplateRepo:        NewOrganizationTemplateRepo(db),
		OrganizationSettingRepo:         NewOrganizationSettingRepo(db),
		BorrowerPropertyTransactionRepo: NewBorrowerPropertyTransactionRepo(db),
		OperationRepo:                   NewOperationRepo(db),
		DocumentRepo:                    NewDocumentRepo(db),
		EntityDocumentRepo:              NewEntityDocumentRepo(db),
		DocumentTemplateRepo:            NewDocumentTemplateRepo(db),
		DriveAuthorizationRepo:          NewDriveAuthorizationRepo(db),
		DriveSyncJobRepo:                NewDriveSyncJobRepo(db),
		PricingPlanRepo:                 NewPricingPlanRepo(db),
		PricingPlanPolicyRepo:           NewPricingPlanPolicyRepo(db),
		PricingPlanAndResourceRepo:      NewPricingPlanAndResourceRepo(db),
		OrganizationSubscriptionRepo:    NewOrganizationSubscriptionRepo(db),
		StripeEventRepo:                 NewStripeEventRepo(db),
		PaymentTransactionRepo:          NewPaymentTransactionRepo(db),
	}
}
