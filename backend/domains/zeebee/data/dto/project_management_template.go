package dto

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementTemplateDto struct {
	Id              int64
	TemplateType    constants.ProjectManagementTemplateType
	OrganizationId  int64
	Name            string
	Description     string
	Status          string
	IsDefault       bool
	SortOrder       int
	ProjectType     constants.ProjectManagementProjectType
	StartAt         *sql.NullTime
	EndAt           *sql.NullTime
	PeriodValue     int
	PeriodUnit      constants.ProjectManagementPeriodUnit
	SourceProjectId int64
}

type ProjectManagementTemplateCreateDto struct {
	TemplateType    constants.ProjectManagementTemplateType
	OrganizationId  int64
	Name            string
	Description     string
	Status          string
	IsDefault       bool
	SortOrder       int
	ProjectType     constants.ProjectManagementProjectType
	StartAt         *sql.NullTime
	EndAt           *sql.NullTime
	PeriodValue     int
	PeriodUnit      constants.ProjectManagementPeriodUnit
	SourceProjectId int64
	CopyTasks       bool
	CopyAssignees   bool
}

type ProjectManagementTemplateUpdateDto struct {
	Name        string
	Description string
	Status      string
	IsDefault   bool
	SortOrder   int
	ProjectType constants.ProjectManagementProjectType
	StartAt     *sql.NullTime
	EndAt       *sql.NullTime
	PeriodValue int
	PeriodUnit  constants.ProjectManagementPeriodUnit
	Stages      []ProjectManagementTemplateStageDto
}

type ProjectManagementTemplateStageDto struct {
	Id          int64
	Name        string
	Description string
	Phase       constants.ProjectManagementStagePhase
	Color       string
	SortOrder   int
}
