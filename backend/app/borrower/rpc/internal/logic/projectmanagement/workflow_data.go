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
	return copierx.Trans[model.ProjectManagementWorkflow, borrower.ProjectManagementWorkflowEntity](workflow)
}
