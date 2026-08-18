package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementProject struct {
	Base
	OrganizationId int64                             `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	PropertyId     int64                             `gorm:"column:property_id;not null;default:0" json:"property_id"`
	Name           string                            `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description    string                            `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	CoverUrl       string                            `gorm:"column:cover_url;type:varchar(255);not null;default:''" json:"cover_url"`
	OwnerId        int64                             `gorm:"column:owner_id;not null;default:0" json:"owner_id"`
	Members        []*ProjectManagementProjectMember `gorm:"foreignKey:ProjectId" json:"members"`
}

func (ProjectManagementProject) TableName() string { return "project_management_projects" }

type ProjectManagementProjectMember struct {
	Base
	ProjectId          int64                                 `gorm:"column:project_id;not null" json:"project_id"`
	OrganizationUserId int64                                 `gorm:"column:organization_user_id;not null" json:"organization_user_id"`
	Role               constants.ProjectManagementMemberRole `gorm:"column:role;type:varchar(32);not null" json:"role"`
}

func (ProjectManagementProjectMember) TableName() string {
	return "project_management_project_members"
}
