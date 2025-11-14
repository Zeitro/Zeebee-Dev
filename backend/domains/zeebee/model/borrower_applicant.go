package model

import (
	"database/sql"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type BorrowerApplicant struct {
	Base
	BorrowerId       int64                   `gorm:"column:borrower_id;default:0" json:"borrower_id"`
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
}

func (m BorrowerApplicant) TableName() string {
	return "borrower_applicants"
}
