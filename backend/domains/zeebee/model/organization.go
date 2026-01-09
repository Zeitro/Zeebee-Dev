package model

import "zeebee-backend/domains/zeebee/constants"

type OrganizationMetadata struct {
	Lender *OrganizationLenderMetadata `json:"lender,omitempty"`
}

type OrganizationLenderMetadata struct {
	LoanAmountRange RangeFloat64 `json:"loan_amount_range"`
	LoanTermRange   RangeInt     `json:"loan_term_range"`
	LoanTateRange   RangeFloat64 `json:"loan_tate_range"`
	LoanPointRange  RangeFloat64 `json:"loan_point_range"`
}

type Organization struct {
	Base
	Name         string                          `gorm:"column:name;default:''" json:"name"`
	Status       constants.OrganizationStatus    `gorm:"column:status;default:''" json:"status"`
	OwnerId      int64                           `gorm:"column:owner_id;default:0" json:"owner_id"`
	Address      Address                         `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	ContactEmail string                          `gorm:"column:contact_email;default:''" json:"contact_email"`
	ContactPhone string                          `gorm:"column:contact_phone;type:varchar(20);default:''" json:"contact_phone"`
	ContactName  string                          `gorm:"column:contact_name;type:varchar(100);default:''" json:"contact_name"`
	OwnerType    constants.OrganizationOwnerType `gorm:"column:owner_type;default:''" json:"owner_type"`
	Type         constants.UserType              `gorm:"column:type;default:''" json:"type"`
	Logo         string                          `gorm:"column:logo;type:varchar(255);default:''" json:"logo"`
	Website      string                          `gorm:"column:website;type:varchar(255);default:''" json:"website"`
	Tagline      string                          `gorm:"column:tagline;type:varchar(255);default:''" json:"tagline"`
	Metadata     *OrganizationMetadata           `gorm:"column:metadata;type:jsonb;default:'{}';serializer:json" json:"metadata"`

	//  relationships
	Users              []*User              `json:"users" gorm:"many2many:organization_and_users;joinForeignKey:organization_id;joinReferences:user_id"`
	OrganizationEmails []*OrganizationEmail `json:"email_authorizations" gorm:"foreignKey:organization_id"`
	Owner              *User                `json:"owner" gorm:"foreignKey:OwnerId;references:Id"`
}

func (m Organization) TableName() string {
	return "organizations"
}
func (m *Organization) RemoveRelationships() {
	m.Users = nil
	m.OrganizationEmails = nil
	m.Owner = nil
}
