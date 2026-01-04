package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type User struct {
	Base
	Type                  constants.UserType   `json:"type"`
	Email                 string               `json:"email"`
	PasswordHash          string               `json:"-"`
	Status                constants.UserStatus `json:"status"`
	FirstName             string               `json:"first_name" gorm:"column:first_name;varchar(100);not null;default:''"`
	MidName               string               `json:"mid_name" gorm:"column:mid_name;varchar(100);not null;default:''"`
	LastName              string               `json:"last_name" gorm:"column:last_name;varchar(100);not null;default:''"`
	FullName              string               `json:"full_name" gorm:"column:full_name;varchar(255);not null;default:''"`
	Phone                 string               `json:"phone" gorm:"column:phone;varchar(20);not null;default:''"`
	DefaultOrganizationId int64                `json:"default_organization_id" gorm:"column:default_organization_id;default:0"`
	// relationships
	Organizations     []*Organization     `json:"organizations" gorm:"many2many:organization_users;joinForeignKey:user_id;joinReferences:organization_id;"`
	OrganizationUsers []*OrganizationUser `json:"organization_users" gorm:"foreignKey:UserId;references:Id"`
}

func (m User) TableName() string {
	return "users"
}

func (m *User) RemoveRelationships() {
	m.Organizations = nil
	m.OrganizationUsers = nil
	// m.Roles = nil
}
