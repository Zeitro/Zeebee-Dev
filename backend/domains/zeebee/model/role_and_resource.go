package model

import "zeebee-backend/domains/zeebee/constants"

type RoleAndResource struct {
	Id         int64               `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	RoleId     int64               `gorm:"column:role_id;default:0" json:"role_id"`
	ResourceId int64               `gorm:"column:resource_id;default:0" json:"resource_id"`
	DataScope  constants.DataScope `gorm:"column:data_scope;default:'All'" json:"data_scope"`
}

func (m RoleAndResource) TableName() string {
	return "role_and_resources"
}
