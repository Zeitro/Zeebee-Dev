package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationUser struct {
	Base
	OrganizationId int64                            `gorm:"column:organization_id;default:0" json:"organization_id"`
	UserId         int64                            `gorm:"column:user_id;default:0" json:"user_id"`
	Status         constants.OrganizationUserStatus `gorm:"column:status;default:'Not Activated'" json:"status"`
	FirstName      string                           `gorm:"column:first_name;default:''" json:"first_name"`
	MidName        string                           `gorm:"column:mid_name;default:''" json:"mid_name"`
	LastName       string                           `gorm:"column:last_name;default:''" json:"last_name"`
	FullName       string                           `gorm:"column:full_name;default:''" json:"full_name"`
	Phone          string                           `gorm:"column:phone;default:''" json:"phone"`
	Remark         string                           `gorm:"column:remark;default:''" json:"remark"`

	// relations
	User         *User         `gorm:"foreignKey:UserId;references:Id" json:"user"`
	Organization *Organization `gorm:"foreignKey:OrganizationId;references:Id" json:"organization"`
	Roles        []*Role       `json:"roles" gorm:"many2many:organization_user_and_roles;joinForeignKey:organization_user_id;joinReferences:role_id;"`
}

func (m OrganizationUser) TableName() string {
	return "organization_users"
}

func (m *OrganizationUser) RemoveRelationships() {
	m.User = nil
	m.Organization = nil
	m.Roles = nil
}
