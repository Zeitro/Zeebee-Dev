package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type BorrowerApplicant struct {
	Base
	CreatorId        int64                   `gorm:"column:creator_id;type:int4;default:0;not null" json:"creator_id"`
	CreatorName      string                  `gorm:"column:creator_name;type:varchar(100);default:'';not null" json:"creator_name"`
	OrganizationId   int64                   `gorm:"column:organization_id;type:int4;default:0;not null" json:"organization_id"`
	FirstName        string                  `gorm:"column:first_name;type:varchar(100);default:'';not null" json:"first_name"`
	MidName          string                  `gorm:"column:mid_name;type:varchar(100);default:'';not null" json:"mid_name"`
	LastName         string                  `gorm:"column:last_name;type:varchar(100);default:'';not null" json:"last_name"`
	FullName         string                  `gorm:"column:full_name;type:varchar(200);default:'';not null" json:"full_name"`
	Ssn              string                  `gorm:"column:ssn;type:varchar(20);default:'';not null" json:"ssn"`
	Dob              *sql.NullTime           `gorm:"column:dob;type:jsonb;default:'{}'" json:"dob"`
	MaritalStatus    constants.MaritalStatus `gorm:"column:marital_status;type:varchar(20);default:'';not null" json:"marital_status"`
	Email            string                  `gorm:"column:email;type:varchar(100);default:'';not null" json:"email"`
	Phone            string                  `gorm:"column:phone;type:varchar(20);default:'';not null" json:"phone"`
	EstCreditScore   int32                   `gorm:"column:est_credit_score;type:int2;default:0;not null" json:"est_credit_score"`
	DriverLicenseNum string                  `gorm:"column:driver_license_num;type:varchar(30);default:'';not null" json:"driver_license_num"`
	Address          Address                 `gorm:"column:address;type:jsonb;serializer:json" json:"address"`

	//  relationships
	Creator      *User         `json:"creator" gorm:"foreignKey:CreatorId;references:Id"`
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationId;references:Id"`
}

func (m BorrowerApplicant) TableName() string {
	return "borrower_applicants"
}

func (m *BorrowerApplicant) RemoveRelationships() {
	m.Creator = nil
	m.Organization = nil
}
