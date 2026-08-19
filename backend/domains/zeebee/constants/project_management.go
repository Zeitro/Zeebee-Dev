package constants

type ProjectManagementWorkflowKind string
type ProjectManagementStagePhase string
type ProjectManagementMemberRole string
type ProjectManagementDependencyType string
type ProjectManagementDependencyRule string

const (
	ProjectManagementWorkflowKindOrganizationTemplate ProjectManagementWorkflowKind = "OrganizationTemplate"
	ProjectManagementWorkflowKindProject              ProjectManagementWorkflowKind = "Project"

	ProjectManagementStagePhasePending    ProjectManagementStagePhase = "Pending"
	ProjectManagementStagePhaseInProgress ProjectManagementStagePhase = "InProgress"
	ProjectManagementStagePhaseFinished   ProjectManagementStagePhase = "Finished"

	ProjectManagementMemberRoleProjectOwner ProjectManagementMemberRole = "ProjectOwner"
	ProjectManagementMemberRoleProjectAdmin ProjectManagementMemberRole = "ProjectAdmin"
	ProjectManagementMemberRoleTaskExecutor ProjectManagementMemberRole = "TaskExecutor"
	ProjectManagementMemberRoleTaskCreator  ProjectManagementMemberRole = "TaskCreator"

	ProjectManagementDependencyTypePreTask  ProjectManagementDependencyType = "PreTask"
	ProjectManagementDependencyTypePostTask ProjectManagementDependencyType = "PostTask"

	ProjectManagementDependencyRuleFS ProjectManagementDependencyRule = "FS"
	ProjectManagementDependencyRuleFF ProjectManagementDependencyRule = "FF"
	ProjectManagementDependencyRuleSS ProjectManagementDependencyRule = "SS"
	ProjectManagementDependencyRuleSF ProjectManagementDependencyRule = "SF"

	ProjectManagementEntityTypeProperty = "Property"

	OperationEntityTypeProjectManagementOrganization OperationEntityType = "ProjectManagementOrganization"
	OperationEntityTypeProjectManagementProject      OperationEntityType = "ProjectManagementProject"
	FileEntityTypeProjectManagementTask              FileEntityType      = "ProjectManagementTask"

	OperationActionProjectManagementTaskTypeCreate   OperationAction = "Project task type created"
	OperationActionProjectManagementTaskTypeUpdate   OperationAction = "Project task type updated"
	OperationActionProjectManagementTaskTypeDelete   OperationAction = "Project task type deleted"
	OperationActionProjectManagementWorkflowCreate   OperationAction = "Project workflow created"
	OperationActionProjectManagementWorkflowUpdate   OperationAction = "Project workflow updated"
	OperationActionProjectManagementWorkflowDelete   OperationAction = "Project workflow deleted"
	OperationActionProjectManagementProjectCreate    OperationAction = "Project created"
	OperationActionProjectManagementProjectUpdate    OperationAction = "Project updated"
	OperationActionProjectManagementProjectDelete    OperationAction = "Project deleted"
	OperationActionProjectManagementMemberAdd        OperationAction = "Project member added"
	OperationActionProjectManagementMemberUpdate     OperationAction = "Project member role updated"
	OperationActionProjectManagementMemberRemove     OperationAction = "Project member removed"
	OperationActionProjectManagementOwnerTransfer    OperationAction = "Project owner transferred"
	OperationActionProjectManagementTaskCreate       OperationAction = "Project task created"
	OperationActionProjectManagementTaskUpdate       OperationAction = "Project task updated"
	OperationActionProjectManagementTaskDelete       OperationAction = "Project task deleted"
	OperationActionProjectManagementTaskAssignees    OperationAction = "Project task assignees updated"
	OperationActionProjectManagementAttachmentAdd    OperationAction = "Project task attachment added"
	OperationActionProjectManagementAttachmentDelete OperationAction = "Project task attachment deleted"
	OperationActionProjectManagementDependencyCreate OperationAction = "Project task dependency created"
	OperationActionProjectManagementDependencyUpdate OperationAction = "Project task dependency updated"
	OperationActionProjectManagementDependencyDelete OperationAction = "Project task dependency deleted"
	OperationActionProjectManagementTaskTransition   OperationAction = "Project task transitioned"
)
