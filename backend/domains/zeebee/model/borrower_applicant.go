package model

import (
	"database/sql"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type BorrowerApplicant struct {
	Base
	UserId           int64                   `gorm:"column:user_id;default:0" json:"user_id"`
	OrganizationId   int64                   `gorm:"column:organization_id;default:0" json:"organization_id"`
	FirstName        string                  `gorm:"column:first_name;default:''" json:"first_name"`
	MidName          string                  `gorm:"column:mid_name;default:''" json:"mid_name"`
	LastName         string                  `gorm:"column:last_name;default:''" json:"last_name"`
	FullName         string                  `gorm:"column:full_name;default:''" json:"full_name"`
	Ssn              string                  `gorm:"column:ssn;default:''" json:"ssn"`
	Dob              *sql.NullTime           `gorm:"column:dob" json:"dob"`
	MaritalStatus    constants.MaritalStatus `gorm:"column:marital_status;default:''" json:"marital_status"`
	Email            string                  `gorm:"column:email;default:''" json:"email"`
	Phone            string                  `gorm:"column:phone;default:''" json:"phone"`
	EstCreditScore   int32                   `gorm:"column:est_credit_score;default:0" json:"est_credit_score"`
	DriverLicenseNum string                  `gorm:"column:driver_license_num;default:''" json:"driver_license_num"`
	Address          types.Address           `gorm:"column:address;type:jsonb;serializer:json" json:"address"`

	//  relationships
	User         *User         `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationId;references:Id"`
}

func (m BorrowerApplicant) TableName() string {
	return "borrower_applicants"
}

func (m *BorrowerApplicant) RemoveRelationships() {
	m.User = nil
	m.Organization = nil
}
