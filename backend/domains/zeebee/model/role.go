package model

import "zeebee-backend/common/constants"

type Role struct {
	Base
	Name string             `gorm:"column:name;type:varchar(100);default:''" json:"name"`
	Type constants.UserType `gorm:"column:type;type:varchar(20);default:'Borrower'" json:"type"`
	// DataScope   constants.DataScope `gorm:"column:data_scope;type:varchar(10);default:'All'" json:"data_scope"`
	Description string `gorm:"column:description;type:varchar(255);default:''" json:"description"`
	//  relationships
	Resources []*Resource `json:"resources" gorm:"many2many:role_and_resources;joinForeignKey:role_id;joinReferences:resource_id"`
	Users     []*User     `json:"users" gorm:"many2many:user_and_roles;joinForeignKey:role_id;joinReferences:user_id"`
}

func (Role) TableName() string {
	return "roles"
}
func (m *Role) RemoveRelationships() {
	m.Resources = nil
	m.Users = nil
}
