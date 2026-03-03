package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type OrganizationEntityDto struct {
	Id                   int64                           `json:"id"`
	Name                 string                          `json:"name"`
	Status               constants.OrganizationStatus    `json:"status"`
	OwnerId              int64                           `json:"owner_id"`
	Address              data.Address                    `json:"address"`
	ContactEmail         string                          `json:"contact_email"`
	ContactPhone         string                          `json:"contact_phone"`
	ContactName          string                          `json:"contact_name"`
	OwnerType            constants.OrganizationOwnerType `json:"owner_type"`
	Type                 constants.UserType              `json:"type"`
	OrganizationEmails   []*OrganizationEmailEntityDto   `json:"organization_emails"`
	OrganizationSettings []*OrganizationSettingEntityDto `json:"organization_settings"`
	Logo                 string                          `json:"logo"`
	Website              string                          `json:"website"`
	Tagline              string                          `json:"tagline"`
	Metadata             *OrganizationMetadata           `json:"metadata"`
	Owner                *UserEntityDto                  `json:"owner"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrganizationPageEntityDto struct {
	Id             int64                           `json:"id"`
	Name           string                          `json:"name"`
	Status         constants.OrganizationStatus    `json:"status"`
	Address        data.Address                    `json:"address"`
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

type OrganizationOptionDto struct {
	Id           int64              `json:"id"`
	Name         string             `json:"name"`
	Type         constants.UserType `json:"type"`
	ContactEmail string             `json:"contact_email"`
}

type OrganizationPageRequestDto struct {
	BasePageRequestDto
	Sorter BasePageSorter
	Filter OrganizationPageFilter
}
type OrganizationPageFilter struct {
	Keyword          string                       `json:"keyword"`
	OrganizationType constants.UserType           `json:"organization_type"`
	Status           constants.OrganizationStatus `json:"status"`
}
type OrganizationRegisterRequestDto struct {
	// Id           int64                           `json:"id"`
	Name         string                          `json:"name"`
	Status       constants.OrganizationStatus    `json:"status"`
	OwnerType    constants.OrganizationOwnerType `json:"owner_type"`
	Type         constants.UserType              `json:"type"`
	ContactEmail string                          `json:"contact_email"`
	ContactPhone string                          `json:"contact_phone"`
	ContactName  string                          `json:"contact_name"`
	Address      *AddressDto                     `json:"address"`
	Logo         string                          `json:"logo"`
	Website      string                          `json:"website"`
	Tagline      string                          `json:"tagline"`
	Metadata     *OrganizationMetadata           `json:"metadata,omitempty"`
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
