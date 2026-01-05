package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type ApiWhiteList struct {
	Base
	Name        string                 `gorm:"column:name;type:varchar(100);default:'';not null" json:"name"`
	Type        constants.UserType     `gorm:"column:type;type:varchar(20);default:'Borrower';not null" json:"type"`
	Path        string                 `gorm:"column:path;type:varchar(100);default:'';not null" json:"path"`
	Method      string                 `gorm:"column:method;type:varchar(20);default:'GET';not null" json:"method"`
	Description string                 `gorm:"column:description;type:varchar(255);default:'';not null" json:"description"`
	Status      constants.CommonStatus `gorm:"column:status;type:varchar(20);default:'Active';not null" json:"status"`
}

func (ApiWhiteList) TableName() string {
	return "api_white_lists"
}

func (m *ApiWhiteList) RemoveRelationships() {

}
