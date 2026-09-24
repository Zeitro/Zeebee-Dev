package constants

import "zeebee-backend/common/enums"

type ProjectManagementWorkflowKind string

const ProjectManagementWorkflowKindName = "project_management_workflow_kind"

const (
	ProjectManagementWorkflowKindOrganizationTemplate ProjectManagementWorkflowKind = "OrganizationTemplate"
	ProjectManagementWorkflowKindProject              ProjectManagementWorkflowKind = "Project"
)

func (c ProjectManagementWorkflowKind) Validate() error {
	return validateConstant(c, []ProjectManagementWorkflowKind{
		ProjectManagementWorkflowKindOrganizationTemplate,
		ProjectManagementWorkflowKindProject,
	}, ProjectManagementWorkflowKindName)
}

type ProjectManagementOwnerType string

const ProjectManagementOwnerTypeName = "project_management_owner_type"

const (
	ProjectManagementOwnerTypeProject  ProjectManagementOwnerType = "Project"
	ProjectManagementOwnerTypeTemplate ProjectManagementOwnerType = "Template"
	ProjectManagementOwnerTypeOrganization ProjectManagementOwnerType = "Organization"
)

func (c ProjectManagementOwnerType) Validate() error {
	return validateConstant(c, []ProjectManagementOwnerType{
		ProjectManagementOwnerTypeProject,
		ProjectManagementOwnerTypeTemplate,
		ProjectManagementOwnerTypeOrganization,
	}, ProjectManagementOwnerTypeName)
}

type ProjectManagementTemplateType string

const ProjectManagementTemplateTypeName = "project_management_template_type"

const (
	ProjectManagementTemplateTypePublic       ProjectManagementTemplateType = "Public"
	ProjectManagementTemplateTypeOrganization ProjectManagementTemplateType = "Organization"
)

func (c ProjectManagementTemplateType) Validate() error {
	return validateConstant(c, []ProjectManagementTemplateType{
		ProjectManagementTemplateTypePublic,
		ProjectManagementTemplateTypeOrganization,
	}, ProjectManagementTemplateTypeName)
}

type ProjectManagementProjectType string

const ProjectManagementProjectTypeName = "project_management_project_type"

const (
	ProjectManagementProjectTypeOrdinary ProjectManagementProjectType = "Ordinary"
	ProjectManagementProjectTypePlanned  ProjectManagementProjectType = "Planned"
)

func (c ProjectManagementProjectType) Validate() error {
	return validateConstant(c, []ProjectManagementProjectType{
		ProjectManagementProjectTypeOrdinary,
		ProjectManagementProjectTypePlanned,
	}, ProjectManagementProjectTypeName)
}

type ProjectManagementPeriodUnit string

const ProjectManagementPeriodUnitName = "project_management_period_unit"

const (
	ProjectManagementPeriodUnitDay   ProjectManagementPeriodUnit = "Day"
	ProjectManagementPeriodUnitWeek  ProjectManagementPeriodUnit = "Week"
	ProjectManagementPeriodUnitMonth ProjectManagementPeriodUnit = "Month"
)

func (c ProjectManagementPeriodUnit) Validate() error {
	return validateConstant(c, []ProjectManagementPeriodUnit{
		ProjectManagementPeriodUnitDay,
		ProjectManagementPeriodUnitWeek,
		ProjectManagementPeriodUnitMonth,
	}, ProjectManagementPeriodUnitName)
}

type ProjectManagementStagePhase string

const ProjectManagementStagePhaseName = "project_management_stage_phase"

const (
	ProjectManagementStagePhasePending    ProjectManagementStagePhase = "Pending"
	ProjectManagementStagePhaseInProgress ProjectManagementStagePhase = "InProgress"
	ProjectManagementStagePhaseFinished   ProjectManagementStagePhase = "Finished"
)

func (c ProjectManagementStagePhase) Validate() error {
	return validateConstant(c, []ProjectManagementStagePhase{
		ProjectManagementStagePhasePending,
		ProjectManagementStagePhaseInProgress,
		ProjectManagementStagePhaseFinished,
	}, ProjectManagementStagePhaseName)
}

type ProjectManagementMemberRole string

const ProjectManagementMemberRoleName = "project_management_member_role"

const (
	ProjectManagementMemberRoleProjectOwner ProjectManagementMemberRole = "ProjectOwner"
	ProjectManagementMemberRoleProjectAdmin ProjectManagementMemberRole = "ProjectAdmin"
	ProjectManagementMemberRoleTaskExecutor ProjectManagementMemberRole = "TaskExecutor"
	ProjectManagementMemberRoleTaskCreator  ProjectManagementMemberRole = "TaskCreator"
)

func (c ProjectManagementMemberRole) Validate() error {
	return validateConstant(c, []ProjectManagementMemberRole{
		ProjectManagementMemberRoleProjectOwner,
		ProjectManagementMemberRoleProjectAdmin,
		ProjectManagementMemberRoleTaskExecutor,
		ProjectManagementMemberRoleTaskCreator,
	}, ProjectManagementMemberRoleName)
}

type ProjectManagementMemberCandidatePurpose string

const ProjectManagementMemberCandidatePurposeName = "project_management_member_candidate_purpose"

const (
	ProjectManagementMemberCandidatePurposeAddMember     ProjectManagementMemberCandidatePurpose = "AddMember"
	ProjectManagementMemberCandidatePurposeTransferOwner ProjectManagementMemberCandidatePurpose = "TransferOwner"
	ProjectManagementMemberCandidatePurposeProjectMember ProjectManagementMemberCandidatePurpose = "ProjectMember"
)

func (c ProjectManagementMemberCandidatePurpose) Validate() error {
	return validateConstant(c, []ProjectManagementMemberCandidatePurpose{
		ProjectManagementMemberCandidatePurposeAddMember,
		ProjectManagementMemberCandidatePurposeTransferOwner,
		ProjectManagementMemberCandidatePurposeProjectMember,
	}, ProjectManagementMemberCandidatePurposeName)
}

type ProjectManagementDependencyType string

const ProjectManagementDependencyTypeName = "project_management_dependency_type"

const (
	ProjectManagementDependencyTypePreTask  ProjectManagementDependencyType = "PreTask"
	ProjectManagementDependencyTypePostTask ProjectManagementDependencyType = "PostTask"
)

func (c ProjectManagementDependencyType) Validate() error {
	return validateConstant(c, []ProjectManagementDependencyType{
		ProjectManagementDependencyTypePreTask,
		ProjectManagementDependencyTypePostTask,
	}, ProjectManagementDependencyTypeName)
}

type ProjectManagementDependencyRule string

const ProjectManagementDependencyRuleName = "project_management_dependency_rule"

const (
	ProjectManagementDependencyRuleFS ProjectManagementDependencyRule = "FS"
	ProjectManagementDependencyRuleFF ProjectManagementDependencyRule = "FF"
	ProjectManagementDependencyRuleSS ProjectManagementDependencyRule = "SS"
	ProjectManagementDependencyRuleSF ProjectManagementDependencyRule = "SF"
)

func (c ProjectManagementDependencyRule) Validate() error {
	return validateConstant(c, []ProjectManagementDependencyRule{
		ProjectManagementDependencyRuleFS,
		ProjectManagementDependencyRuleFF,
		ProjectManagementDependencyRuleSS,
		ProjectManagementDependencyRuleSF,
	}, ProjectManagementDependencyRuleName)
}

type ProjectManagementProjectEntityType string

const ProjectManagementProjectEntityTypeName = "project_management_project_entity_type"

const (
	ProjectManagementEntityTypeNone             ProjectManagementProjectEntityType = ""
	ProjectManagementEntityTypeBorrowerProperty ProjectManagementProjectEntityType = "BorrowerProperty"
)

func (c ProjectManagementProjectEntityType) Validate() error {
	return validateConstant(c, []ProjectManagementProjectEntityType{
		ProjectManagementEntityTypeBorrowerProperty,
	}, ProjectManagementProjectEntityTypeName)
}

// ProjectManagementPriority task priority.
type ProjectManagementPriority string

const ProjectManagementPriorityName = "project_management_priority"

const (
	ProjectManagementPriorityLow      ProjectManagementPriority = "Low"
	ProjectManagementPriorityMedium   ProjectManagementPriority = "Medium"
	ProjectManagementPriorityHigh     ProjectManagementPriority = "High"
	ProjectManagementPriorityCritical ProjectManagementPriority = "Critical"
)

func (c ProjectManagementPriority) Validate() error {
	return validateConstant(c, []ProjectManagementPriority{
		ProjectManagementPriorityLow,
		ProjectManagementPriorityMedium,
		ProjectManagementPriorityHigh,
		ProjectManagementPriorityCritical,
	}, ProjectManagementPriorityName)
}

func init() {
	enums.Register(ProjectManagementWorkflowKindName, []enums.EnumOpt{
		{Label: string(ProjectManagementWorkflowKindOrganizationTemplate), Value: string(ProjectManagementWorkflowKindOrganizationTemplate)},
		{Label: string(ProjectManagementWorkflowKindProject), Value: string(ProjectManagementWorkflowKindProject)},
	})
	enums.Register(ProjectManagementOwnerTypeName, []enums.EnumOpt{
		{Label: string(ProjectManagementOwnerTypeProject), Value: string(ProjectManagementOwnerTypeProject)},
		{Label: string(ProjectManagementOwnerTypeTemplate), Value: string(ProjectManagementOwnerTypeTemplate)},
		{Label: string(ProjectManagementOwnerTypeOrganization), Value: string(ProjectManagementOwnerTypeOrganization)},
	})
	enums.Register(ProjectManagementTemplateTypeName, []enums.EnumOpt{
		{Label: string(ProjectManagementTemplateTypePublic), Value: string(ProjectManagementTemplateTypePublic)},
		{Label: string(ProjectManagementTemplateTypeOrganization), Value: string(ProjectManagementTemplateTypeOrganization)},
	})
	enums.Register(ProjectManagementProjectTypeName, []enums.EnumOpt{
		{Label: string(ProjectManagementProjectTypeOrdinary), Value: string(ProjectManagementProjectTypeOrdinary)},
		{Label: string(ProjectManagementProjectTypePlanned), Value: string(ProjectManagementProjectTypePlanned)},
	})
	enums.Register(ProjectManagementPeriodUnitName, []enums.EnumOpt{
		{Label: string(ProjectManagementPeriodUnitDay), Value: string(ProjectManagementPeriodUnitDay)},
		{Label: string(ProjectManagementPeriodUnitWeek), Value: string(ProjectManagementPeriodUnitWeek)},
		{Label: string(ProjectManagementPeriodUnitMonth), Value: string(ProjectManagementPeriodUnitMonth)},
	})
	enums.Register(ProjectManagementStagePhaseName, []enums.EnumOpt{
		{Label: string(ProjectManagementStagePhasePending), Value: string(ProjectManagementStagePhasePending)},
		{Label: string(ProjectManagementStagePhaseInProgress), Value: string(ProjectManagementStagePhaseInProgress)},
		{Label: string(ProjectManagementStagePhaseFinished), Value: string(ProjectManagementStagePhaseFinished)},
	})
	enums.Register(ProjectManagementMemberRoleName, []enums.EnumOpt{
		{Label: string(ProjectManagementMemberRoleProjectOwner), Value: string(ProjectManagementMemberRoleProjectOwner)},
		{Label: string(ProjectManagementMemberRoleProjectAdmin), Value: string(ProjectManagementMemberRoleProjectAdmin)},
		{Label: string(ProjectManagementMemberRoleTaskExecutor), Value: string(ProjectManagementMemberRoleTaskExecutor)},
		{Label: string(ProjectManagementMemberRoleTaskCreator), Value: string(ProjectManagementMemberRoleTaskCreator)},
	})
	enums.Register(ProjectManagementMemberCandidatePurposeName, []enums.EnumOpt{
		{Label: string(ProjectManagementMemberCandidatePurposeAddMember), Value: string(ProjectManagementMemberCandidatePurposeAddMember)},
		{Label: string(ProjectManagementMemberCandidatePurposeTransferOwner), Value: string(ProjectManagementMemberCandidatePurposeTransferOwner)},
		{Label: string(ProjectManagementMemberCandidatePurposeProjectMember), Value: string(ProjectManagementMemberCandidatePurposeProjectMember)},
	})
	enums.Register(ProjectManagementDependencyTypeName, []enums.EnumOpt{
		{Label: string(ProjectManagementDependencyTypePreTask), Value: string(ProjectManagementDependencyTypePreTask)},
		{Label: string(ProjectManagementDependencyTypePostTask), Value: string(ProjectManagementDependencyTypePostTask)},
	})
	enums.Register(ProjectManagementDependencyRuleName, []enums.EnumOpt{
		{Label: string(ProjectManagementDependencyRuleFS), Value: string(ProjectManagementDependencyRuleFS)},
		{Label: string(ProjectManagementDependencyRuleFF), Value: string(ProjectManagementDependencyRuleFF)},
		{Label: string(ProjectManagementDependencyRuleSS), Value: string(ProjectManagementDependencyRuleSS)},
		{Label: string(ProjectManagementDependencyRuleSF), Value: string(ProjectManagementDependencyRuleSF)},
	})
	enums.Register(ProjectManagementProjectEntityTypeName, []enums.EnumOpt{
		{Label: string(ProjectManagementEntityTypeBorrowerProperty), Value: string(ProjectManagementEntityTypeBorrowerProperty)},
	})
	enums.Register(ProjectManagementPriorityName, []enums.EnumOpt{
		{Label: string(ProjectManagementPriorityLow), Value: string(ProjectManagementPriorityLow)},
		{Label: string(ProjectManagementPriorityMedium), Value: string(ProjectManagementPriorityMedium)},
		{Label: string(ProjectManagementPriorityHigh), Value: string(ProjectManagementPriorityHigh)},
		{Label: string(ProjectManagementPriorityCritical), Value: string(ProjectManagementPriorityCritical)},
	})
}
