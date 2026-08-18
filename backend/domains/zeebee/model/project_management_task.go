package model

import (
	"database/sql"

	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementTaskType struct {
	Base
	OrganizationId int64  `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	Name           string `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description    string `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	Color          string `gorm:"column:color;type:varchar(32);not null;default:''" json:"color"`
	Enabled        bool   `gorm:"column:enabled;not null;default:true" json:"enabled"`
}

func (ProjectManagementTaskType) TableName() string { return "project_management_task_types" }

type ProjectManagementTask struct {
	Base
	ProjectId      int64                            `gorm:"column:project_id;not null" json:"project_id"`
	OrganizationId int64                            `gorm:"column:organization_id;not null" json:"organization_id"`
	ParentTaskId   int64                            `gorm:"column:parent_task_id;not null;default:0" json:"parent_task_id"`
	StageId        int64                            `gorm:"column:stage_id;not null" json:"stage_id"`
	CreatorId      int64                            `gorm:"column:creator_id;not null" json:"creator_id"`
	CreatorName    string                           `gorm:"column:creator_name;type:varchar(100);not null;default:''" json:"creator_name"`
	Title          string                           `gorm:"column:title;type:varchar(255);not null;default:''" json:"title"`
	Description    string                           `gorm:"column:description;type:text;not null;default:''" json:"description"`
	Priority       string                           `gorm:"column:priority;type:varchar(32);not null;default:''" json:"priority"`
	StartAt        *sql.NullTime                    `gorm:"column:start_at;type:timestamptz" json:"start_at"`
	EndAt          *sql.NullTime                    `gorm:"column:end_at;type:timestamptz" json:"end_at"`
	TaskType       datatypes.JSON                   `gorm:"column:task_type;type:jsonb;not null;default:'{}'" json:"task_type"`
	Assignees      []*ProjectManagementTaskAssignee `gorm:"foreignKey:TaskId" json:"assignees"`
	Overdue        bool                             `gorm:"-" json:"overdue"`
	OverdueSeconds int64                            `gorm:"-" json:"overdue_seconds"`
}

func (ProjectManagementTask) TableName() string { return "project_management_tasks" }

type ProjectManagementTaskAssignee struct {
	Base
	TaskId               int64  `gorm:"column:task_id;not null" json:"task_id"`
	OrganizationUserId   int64  `gorm:"column:organization_user_id;not null" json:"organization_user_id"`
	OrganizationUserName string `gorm:"column:organization_user_name;type:varchar(100);not null;default:''" json:"organization_user_name"`
}

func (ProjectManagementTaskAssignee) TableName() string { return "project_management_task_assignees" }

type ProjectManagementTaskDependency struct {
	Base
	SuccessorTaskId   int64                                     `gorm:"column:successor_task_id;not null" json:"successor_task_id"`
	PredecessorTaskId int64                                     `gorm:"column:predecessor_task_id;not null" json:"predecessor_task_id"`
	DependencyType    constants.ProjectManagementDependencyType `gorm:"column:dependency_type;type:varchar(2);not null" json:"dependency_type"`
}

func (ProjectManagementTaskDependency) TableName() string {
	return "project_management_task_dependencies"
}
