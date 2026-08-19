package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/data/dto"
)

type ProjectManagementTask struct {
	Base
	ProjectId      int64                                     `gorm:"column:project_id;not null;default:0" json:"project_id"`
	OrganizationId int64                                     `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	ParentTaskId   int64                                     `gorm:"column:parent_task_id;not null;default:0" json:"parent_task_id"`
	StageId        int64                                     `gorm:"column:stage_id;not null;default:0" json:"stage_id"`
	CreatorId      int64                                     `gorm:"column:creator_id;not null;default:0" json:"creator_id"`
	CreatorName    string                                    `gorm:"column:creator_name;type:varchar(100);not null;default:''" json:"creator_name"`
	Title          string                                    `gorm:"column:title;type:varchar(255);not null;default:''" json:"title"`
	Description    string                                    `gorm:"column:description;type:text;not null;default:''" json:"description"`
	Priority       string                                    `gorm:"column:priority;type:varchar(32);not null;default:''" json:"priority"`
	StartAt        *sql.NullTime                             `gorm:"column:start_at;type:timestamptz" json:"start_at"`
	EndAt          *sql.NullTime                             `gorm:"column:end_at;type:timestamptz" json:"end_at"`
	TaskType       *dto.ProjectManagementTaskTypeSnapshotDto `gorm:"column:task_type;type:jsonb;not null;default:'{}';serializer:json" json:"task_type"`
	// relationships
	Assignees      []*ProjectManagementTaskAssignee `gorm:"foreignKey:TaskId" json:"assignees"`
	// Overdue        bool                             `gorm:"-" json:"overdue"`
	// OverdueSeconds int64                            `gorm:"-" json:"overdue_seconds"`
}

func (ProjectManagementTask) TableName() string { return "project_management_tasks" }

// RemoveRelationships removes relationship data to avoid circular references.
func (m *ProjectManagementTask) RemoveRelationships() {
	m.Assignees = nil
}
