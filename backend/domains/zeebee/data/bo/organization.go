package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationBo struct {
	Id                 int64                           `json:"id"`
	Name               string                          `json:"name"`
	Status             constants.OrganizationStatus    `json:"status"`
	OwnerId            int64                           `json:"owner_id"`
	Address            *AddressBo                      `json:"address"`
	ContactEmail       string                          `json:"contact_email"`
	ContactPhone       string                          `json:"contact_phone"`
	ContactName        string                          `json:"contact_name"`
	OwnerType          constants.OrganizationOwnerType `json:"owner_type"`
	Type               constants.UserType              `json:"type"`
	OrganizationEmails []*OrganizationEmailBo          `json:"organization_emails"`
	Tagline            string                          `json:"tagline"`
	Metadata           *OrganizationMetadata           `json:"metadata"`
	CreatedAt          time.Time                       `json:"created_at"`
	UpdatedAt          time.Time                       `json:"updated_at"`
}

type OrganizationMetadata struct {
	Lender *OrganizationLenderMetadata `json:"lender,omitempty"`
}

type OrganizationLenderMetadata struct {
	LoanAmountRange RangeFloat64 `json:"loan_amount_range"`
	LoanTermRange   RangeInt     `json:"loan_term_range"`
	LoanTateRange   RangeFloat64 `json:"loan_tate_range"`
	LoanPointRange  RangeFloat64 `json:"loan_point_range"`
}
