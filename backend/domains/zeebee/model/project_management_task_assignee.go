package model

type ProjectManagementTaskAssignee struct {
	Base
	TaskId               int64  `gorm:"column:task_id;not null;default:0" json:"task_id"`
	StageId              int64  `gorm:"column:stage_id;not null;default:0;index" json:"stage_id"`
	OrganizationUserId   int64  `gorm:"column:organization_user_id;not null;default:0" json:"organization_user_id"`
	OrganizationUserName string `gorm:"column:organization_user_name;type:varchar(100);not null;default:''" json:"organization_user_name"`
}

func (ProjectManagementTaskAssignee) TableName() string { return "project_management_task_assignees" }
