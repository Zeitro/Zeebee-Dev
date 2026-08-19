package model

type ProjectManagementTaskType struct {
	Base
	OrganizationId int64  `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	Name           string `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description    string `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	Color          string `gorm:"column:color;type:varchar(32);not null;default:''" json:"color"`
	Enabled        bool   `gorm:"column:enabled;not null;default:true" json:"enabled"`
}

func (ProjectManagementTaskType) TableName() string { return "project_management_task_types" }
