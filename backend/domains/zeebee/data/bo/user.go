package bo

import (
	"time"

	"zeebee-backend/common/constants"
)

type UserBo struct {
	Id        int64                `json:"id"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Type      constants.UserType   `json:"type"`
	Email     string               `json:"email"`
	Status    constants.UserStatus `json:"status"`
	FirstName string               `json:"first_name"`
	MidName   string               `json:"mid_name" `
	LastName  string               `json:"last_name" `
	FullName  string               `json:"full_name" `
	Phone     string               `json:"phone"`
}
