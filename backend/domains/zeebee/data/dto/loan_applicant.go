package dto

import (
	"database/sql"
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type LoanApplicantEntityDto struct {
	Id                  int64                               `json:"id"`
	BorrowerApplicantId int64                               `json:"borrower_applicant_id"`
	FirstName           string                              `json:"first_name"`
	LastName            string                              `json:"last_name"`
	IsPrimary           bool                                `json:"is_primary"`
	CreatedAt           time.Time                           `json:"created_at"`
	UpdatedAt           time.Time                           `json:"updated_at"`
	MidName             string                              `json:"mid_name"`
	InvitationStatus    constants.ApplicantInvitationStatus `json:"invitation_status"`
	FullName            string                              `json:"full_name"`
	Ssn                 string                              `json:"ssn"`
	Dob                 *sql.NullTime                       `json:"dob"`
	MaritalStatus       constants.MaritalStatus             `json:"marital_status"`
	Email               string                              `json:"email"`
	EstCreditScore      int32                               `json:"est_credit_score"`
	DriverLicenseNum    string                              `json:"driver_license_num"`
	Address             data.Address                        `json:"address"`
	LoanBorrowerId      int64                               `json:"loan_borrower_id"`
	Phone               string                              `json:"phone"`
	OwnPercent          int32                               `json:"own_percent"`
}
