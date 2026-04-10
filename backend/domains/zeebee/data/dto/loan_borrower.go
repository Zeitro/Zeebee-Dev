package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type LoanBorrowerEntityDto struct {
	Id             int64                   `json:"id"`
	Type           constants.ApplicantType `json:"type"`
	Name           string                  `json:"name"`
	TinType        constants.TINType       `json:"tin_type"`
	Tin            string                  `json:"tin"`
	MailingAddress data.Address            `json:"mailing_address"`
	ContactEmail   string                  `json:"contact_email"`
	ContactPhone   string                  `json:"contact_phone"`
	ContactName    string                  `json:"contact_name"`
	OwnPercent     int32                   `json:"own_percent"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}
