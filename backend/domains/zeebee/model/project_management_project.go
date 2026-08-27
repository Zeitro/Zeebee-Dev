package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementProject struct {
	Base
	OrganizationId int64                                        `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	EntityId       int64                                        `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	EntityType     constants.ProjectManagementProjectEntityType `gorm:"column:entity_type;type:varchar(100);not null;default:''" json:"entity_type"`
	Name           string                                       `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description    string                                       `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	Priority       constants.ProjectManagementPriority          `gorm:"column:priority;type:varchar(32);not null;default:''" json:"priority"`
	CoverUrl       string                                       `gorm:"column:cover_url;type:varchar(255);not null;default:''" json:"cover_url"`
	OwnerId        int64                                        `gorm:"column:owner_id;not null;default:0" json:"owner_id"`
	OwnerName      string                                       `gorm:"column:owner_name;type:varchar(100);not null;default:''" json:"owner_name"`
	// relationships
	Members []*ProjectManagementProjectMember `gorm:"foreignKey:ProjectId" json:"members"`
}

func (ProjectManagementProject) TableName() string { return "project_management_projects" }

// RemoveRelationships removes relationship data to avoid circular references.
func (m *ProjectManagementProject) RemoveRelationships() {
	m.Members = nil
}
