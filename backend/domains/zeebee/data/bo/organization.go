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
	Logo               string                          `json:"logo"`
	Website            string                          `json:"website"`
	Tagline            string                          `json:"tagline"`
	Metadata           *OrganizationMetadata           `json:"metadata"`
	Owner              *UserBo                         `json:"owner"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrganizationPageEntityBo struct {
	Id             int64                           `json:"id"`
	Name           string                          `json:"name"`
	Status         constants.OrganizationStatus    `json:"status"`
	Address        *AddressBo                      `json:"address"`
	ContactEmail   string                          `json:"contact_email"`
	ContactPhone   string                          `json:"contact_phone"`
	ContactName    string                          `json:"contact_name"`
	OwnerType      constants.OrganizationOwnerType `json:"owner_type"`
	Type           constants.UserType              `json:"type"`
	Logo           string                          `json:"logo"`
	Website        string                          `json:"website"`
	Tagline        string                          `json:"tagline"`
	OwnerEmail     string                          `json:"owner_email"`
	OwnerFullName  string                          `json:"owner_full_name"`
	OwnerFirstName string                          `json:"owner_first_name"`
	OwnerMidName   string                          `json:"owner_mid_name"`
	OwnerLastName  string                          `json:"owner_last_name"`
	CreatedAt      time.Time                       `json:"created_at"`
	UpdatedAt      time.Time                       `json:"updated_at"`
}

type OrganizationMetadata struct {
	Lender *OrganizationLenderMetadata `json:"lender,omitempty"`
}

type OrganizationLenderMetadata struct {
	LoanAmountRange RangeFloat64 `json:"loan_amount_range"`
	LoanTermRange   RangeInt     `json:"loan_term_range"`
	LoanRateRange   RangeFloat64 `json:"loan_rate_range"`
	LoanPointRange  RangeFloat64 `json:"loan_point_range"`
}
