package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type LoanThirdServiceEntityBo struct {
	Id int64 `json:"id"`
	// LoanId  int64  `json:"loan_id"`
	Type      constants.LoanThirdServiceType `json:"type"`
	Name      string                         `json:"name"`
	Email     string                         `json:"email"`
	Phone     string                         `json:"phone"`
	Company   string                         `json:"company"`
	CreatedAt time.Time                      `json:"created_at"`
	UpdatedAt time.Time                      `json:"updated_at"`
}
