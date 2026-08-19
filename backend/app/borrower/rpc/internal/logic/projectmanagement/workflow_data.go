package projectmanagementlogic

import (
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func workflowDto(name, description string, enabled, isDefault bool, stages []*borrower.ProjectManagementWorkflowStageEntity, transitions []*borrower.ProjectManagementWorkflowTransitionEntity) (*dto.ProjectManagementWorkflowDto, error) {
	stageData, err := copierx.Trans[[]*borrower.ProjectManagementWorkflowStageEntity, []dto.ProjectManagementWorkflowStageDto](&stages)
	if err != nil {
		return nil, err
	}
	transitionData, err := copierx.Trans[[]*borrower.ProjectManagementWorkflowTransitionEntity, []dto.ProjectManagementWorkflowTransitionDto](&transitions)
	if err != nil {
		return nil, err
	}
	return &dto.ProjectManagementWorkflowDto{Name: name, Description: description, Enabled: enabled, IsDefault: isDefault, Stages: *stageData, Transitions: *transitionData}, nil
}

func workflowEntity(workflow *model.ProjectManagementWorkflow) (*borrower.ProjectManagementWorkflowEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementWorkflow, borrower.ProjectManagementWorkflowEntity](workflow)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(workflow.CreatedAt), formatEntityTime(workflow.UpdatedAt)
	for index, stage := range workflow.Stages {
		entity.Stages[index].CreatedAt = formatEntityTime(stage.CreatedAt)
		entity.Stages[index].UpdatedAt = formatEntityTime(stage.UpdatedAt)
	}
	for index, transition := range workflow.Transitions {
		entity.Transitions[index].CreatedAt = formatEntityTime(transition.CreatedAt)
		entity.Transitions[index].UpdatedAt = formatEntityTime(transition.UpdatedAt)
	}
	return entity, nil
}

func workflowPageEntity(workflow *model.ProjectManagementWorkflow) (*borrower.ProjectManagementWorkflowPageEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementWorkflow, borrower.ProjectManagementWorkflowPageEntity](workflow)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(workflow.CreatedAt), formatEntityTime(workflow.UpdatedAt)
	return entity, nil
}
