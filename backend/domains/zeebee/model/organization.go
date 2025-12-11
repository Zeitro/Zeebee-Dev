package model

import "zeebee-backend/common/constants"

type Organization struct {
	Base
	Name         string                          `gorm:"column:name;default:''" json:"name"`
	Status       constants.OrganizationStatus    `gorm:"column:status;default:''" json:"status"`
	OwnerId      int64                           `gorm:"column:owner_id;default:0" json:"owner_id"`
	Address      Address                         `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	ContactEmail string                          `gorm:"column:contact_email;default:''" json:"contact_email"`
	OwnerType    constants.OrganizationOwnerType `gorm:"column:owner_type;default:''" json:"owner_type"`
	Type         constants.UserType              `gorm:"column:type;default:''" json:"type"`

	//  relationships
	Users []*User `json:"users" gorm:"many2many:organization_and_users;joinForeignKey:organization_id;joinReferences:user_id"`
}

func (m Organization) TableName() string {
	return "organizations"
}
func (m *Organization) RemoveRelationships() {
	m.Users = nil
}
