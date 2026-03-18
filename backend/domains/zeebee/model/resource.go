package model

import "zeebee-backend/domains/zeebee/constants"

type Resource struct {
	Base
	ParentId        int64                  `gorm:"column:parent_id;type:int4;default:0" json:"parent_id"`
	Type            constants.UserType     `gorm:"column:type;type:varchar(20);default:'Borrower'" json:"type"`
	ResourceType    constants.ResourceType `gorm:"column:resource_type;type:varchar(20);default:'Menu'" json:"resource_type"`
	Name            string                 `gorm:"column:name;type:varchar(100);default:''" json:"name"`
	Path            string                 `gorm:"column:path;type:varchar(100);default:''" json:"path"`
	Icon            string                 `gorm:"column:icon;type:varchar(100);default:''" json:"icon"`
	Tag             string                 `gorm:"column:tag;type:varchar(100);default:''" json:"tag"`
	IsShow          bool                   `gorm:"column:is_show;default:false" json:"is_show"`
	Status          constants.CommonStatus `gorm:"column:status;type:varchar(20);default:'Active'" json:"status"`
	SortOrder       int32                  `gorm:"column:sort_order;default:10" json:"sort_order"`
	EnableDataScope bool                   `gorm:"column:enable_data_scope;default:false" json:"enable_data_scope"`
	Description     string                 `gorm:"column:description;type:varchar(255);default:''" json:"description"`
	//  relationships
	Apis  []*Api  `json:"apis" gorm:"many2many:resource_and_apis;joinForeignKey:resource_id;joinReferences:api_id"`
	Roles []*Role `json:"roles" gorm:"many2many:role_and_resources;joinForeignKey:resource_id;joinReferences:role_id"`
}

func (Resource) TableName() string {
	return "resources"
}

func (m *Resource) RemoveRelationships() {
	m.Apis = nil
	m.Roles = nil
}
