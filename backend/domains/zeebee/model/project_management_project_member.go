package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementProjectMember struct {
	Base
	ProjectId          int64                                 `gorm:"column:project_id;not null;default:0" json:"project_id"`
	OrganizationUserId int64                                 `gorm:"column:organization_user_id;not null;default:0" json:"organization_user_id"`
	Role               constants.ProjectManagementMemberRole `gorm:"column:role;type:varchar(32);not null;default:'TaskExecutor'" json:"role"`
}

func (ProjectManagementProjectMember) TableName() string {
	return "project_management_project_members"
}
