package bo

import (
	"time"

	"zeebee-backend/common/constants"
)

type OrganizationUserBo struct {
	Id             int64                            `json:"id"`
	OrganizationId int64                            `json:"organization_id"`
	UserId         int64                            `json:"user_id"`
	Status         constants.OrganizationUserStatus `json:"status"`
	FirstName      string                           `json:"first_name"`
	MidName        string                           `json:"mid_name"`
	LastName       string                           `json:"last_name"`
	FullName       string                           `json:"full_name"`
	Phone          string                           `json:"phone"`
	Remark         string                           `json:"remark"`
	CreatedAt      time.Time                        `json:"created_at"`
	UpdatedAt      time.Time                        `json:"updated_at"`
	Email          string                           `json:"email"`
	IsOwner        bool                             `json:"is_owner"`
	Roles          []*RoleBo                        `json:"roles"`
}
