package model

import "zeebee-backend/common/constants"

type Api struct {
	Base
	ParentId    int64                  `gorm:"column:parent_id;type:int4;default:0;not null" json:"parent_id"`
	Type        constants.UserType     `gorm:"column:type;type:varchar(100);default:'Borrower';not null" json:"type"`
	Name        string                 `gorm:"column:name;type:varchar(100);default:'';not null" json:"name"`
	Path        string                 `gorm:"column:path;type:varchar(100);default:'';not null" json:"path"`
	Method      string                 `gorm:"column:method;type:varchar(20);default:'';not null" json:"method"`
	Description string                 `gorm:"column:description;type:varchar(255);default:'';not null" json:"description"`
	Status      constants.CommonStatus `gorm:"column:status;type:varchar(100);default:'Active';not null" json:"status"`
	SortOrder   int32                  `gorm:"column:sort_order;type:int4;default:10;not null" json:"sort_order"`
	//  relationships
	Resources []*Resource `json:"resources" gorm:"many2many:resource_and_apis;joinForeignKey:api_id;joinReferences:resource_id"`
}

func (Api) TableName() string {
	return "apis"
}

func (m *Api) RemoveRelationships() {
	m.Resources = nil
}
