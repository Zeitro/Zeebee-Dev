package model

import "zeebee-backend/common/constants"

type OrganizationAndUser struct {
	Id               int64                                         `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	OrganizationId   int64                                         `gorm:"column:organization_id;default:0" json:"organization_id"`
	UserId           int64                                         `gorm:"column:user_id;default:0" json:"user_id"`
	Status           constants.OrganizationAndUserStatus           `gorm:"column:status;default:'Active'" json:"status"`
	ActivationStatus constants.OrganizationAndUserActivationStatus `gorm:"column:activation_status;default:'Not Activated'" json:"activation_status"`
}

func (m OrganizationAndUser) TableName() string {
	return "organization_and_users"
}
