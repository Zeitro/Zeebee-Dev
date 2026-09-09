package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementTaskDependency struct {
	Base
	TaskId           int64                                     `gorm:"column:task_id;not null;default:0" json:"task_id"`
	DependencyTaskId int64                                     `gorm:"column:dependency_task_id;not null;default:0" json:"dependency_task_id"`
	DependencyType   constants.ProjectManagementDependencyType `gorm:"column:dependency_type;type:varchar(16);not null;default:'PreTask'" json:"dependency_type"`
	DependencyRule   constants.ProjectManagementDependencyRule `gorm:"column:dependency_rule;type:varchar(2);not null;default:'FS'" json:"dependency_rule"`
	// relationships
	Task           *ProjectManagementTask `gorm:"foreignKey:TaskId;references:Id" json:"task,omitempty"`
	DependencyTask *ProjectManagementTask `gorm:"foreignKey:DependencyTaskId;references:Id" json:"dependency_task,omitempty"`
}

func (ProjectManagementTaskDependency) TableName() string {
	return "project_management_task_dependencies"
}

// RemoveRelationships removes relationship data to avoid circular references.
func (m *ProjectManagementTaskDependency) RemoveRelationships() {
	m.Task = nil
	m.DependencyTask = nil
}
