package model



// ProjectManagementWorkflowStage connects a workflow to an owner-stage.
// Stage definitions belong to a project or template; a workflow only selects them.
type ProjectManagementWorkflowStage struct {
	Id         int64 `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	WorkflowId int64 `gorm:"column:workflow_id;not null;default:0" json:"workflow_id"`
	StageId    int64 `gorm:"column:stage_id;not null;default:0" json:"stage_id"`
}

func (ProjectManagementWorkflowStage) TableName() string {
	return "project_management_workflow_stages"
}

