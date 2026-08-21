package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementWorkflowLogic struct{ *BaseLogic }

type projectManagementWorkflowOperation struct {
	claim      utils.JwtUser
	workflowId int64
	action     constants.OperationAction
}

func NewProjectManagementWorkflowLogic(baseLogic *BaseLogic) *ProjectManagementWorkflowLogic {
	return &ProjectManagementWorkflowLogic{BaseLogic: baseLogic}
}

func ValidateProjectManagementWorkflow(workflow dto.ProjectManagementWorkflowDto) error {
	if strings.TrimSpace(workflow.Name) == "" {
		return errorx.New("workflow name is required")
	}
	stageIDs := make(map[int64]struct{}, len(workflow.Stages))
	sortOrders := make(map[int]struct{}, len(workflow.Stages))
	pendingCount := 0
	for _, stage := range workflow.Stages {
		if stage.Id == 0 || strings.TrimSpace(stage.Name) == "" {
			return errorx.New("workflow stage id and name are required")
		}
		if !isProjectManagementStagePhase(stage.Phase) {
			return errorx.New("invalid workflow stage phase")
		}
		if _, exists := sortOrders[stage.SortOrder]; exists {
			return errorx.New("workflow stage sort order must be unique")
		}
		sortOrders[stage.SortOrder] = struct{}{}
		if _, exists := stageIDs[stage.Id]; exists {
			return errorx.New("workflow stage id must be unique")
		}
		stageIDs[stage.Id] = struct{}{}
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			pendingCount++
		}
	}
	if pendingCount != 1 {
		return errorx.New("workflow must have exactly one pending stage")
	}

	edges := make(map[string]struct{}, len(workflow.Transitions))
	for _, transition := range workflow.Transitions {
		if transition.FromStageId == transition.ToStageId {
			return errorx.New("workflow transition cannot point to itself")
		}
		if _, ok := stageIDs[transition.FromStageId]; !ok {
			return errorx.New("workflow transition source stage is invalid")
		}
		if _, ok := stageIDs[transition.ToStageId]; !ok {
			return errorx.New("workflow transition target stage is invalid")
		}
		edge := fmt.Sprintf("%d:%d", transition.FromStageId, transition.ToStageId)
		if _, exists := edges[edge]; exists {
			return errorx.New("workflow transition must be unique")
		}
		edges[edge] = struct{}{}
		actors := make(map[constants.ProjectManagementMemberRole]struct{}, len(transition.AllowedActorTypes))
		for _, claim := range transition.AllowedActorTypes {
			if !isProjectManagementMemberRole(claim) {
				return errorx.New("invalid workflow transition claim type")
			}
			if _, exists := actors[claim]; exists {
				return errorx.New("workflow transition claim type must be unique")
			}
			actors[claim] = struct{}{}
		}
	}
	return nil
}

func (l *ProjectManagementWorkflowLogic) Page(ctx context.Context, organizationId int64, input dto.ProjectManagementWorkflowPageRequestDto) (gormx.PageResponse[*model.ProjectManagementWorkflow], error) {
	return l.repos.ProjectManagementWorkflowRepo.Page(ctx, organizationId, gormx.PageRequest{Page: input.Page, Limit: input.Limit}, input.Filter.Keyword, input.Filter.Enabled)
}

func (l *ProjectManagementWorkflowLogic) Detail(ctx context.Context, organizationId, id int64) (*model.ProjectManagementWorkflow, error) {
	workflow, err := l.repos.ProjectManagementWorkflowRepo.FindDefinitionByOrganization(ctx, id, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "workflow not found")
	}
	return workflow, nil
}

func (l *ProjectManagementWorkflowLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementWorkflowDto) (workflow *model.ProjectManagementWorkflow, err error) {
	if err = validateProjectManagementWorkflowDetails(input); err != nil {
		return nil, err
	}
	if _, err = l.ListEnabled(ctx, claim.OrganizationId); err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		defaultWorkflow, err := l.repos.ProjectManagementWorkflowRepo.LockDefaultTemplateDefinitionWithTx(ctx, tx, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "default workflow not found")
		}
		workflow = &model.ProjectManagementWorkflow{
			Kind:           constants.ProjectManagementWorkflowKindOrganizationTemplate,
			OrganizationId: claim.OrganizationId,
			Name:           input.Name,
			Description:    input.Description,
			Enabled:        true,
		}
		workflow, err = l.repos.ProjectManagementWorkflowRepo.CreateWithTx(ctx, tx, workflow, false)
		if err != nil {
			return err
		}
		stages, transitions := cloneProjectManagementWorkflowDefinition(defaultWorkflow)
		if err := l.repos.ProjectManagementWorkflowRepo.ReplaceTemplateDefinitionWithTx(ctx, tx, workflow.Id, stages, transitions); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflow.Id, action: constants.OperationActionProjectManagementWorkflowCreate})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim.OrganizationId, workflow.Id)
}

func (l *ProjectManagementWorkflowLogic) Update(ctx context.Context, claim utils.JwtUser, id int64, input dto.ProjectManagementWorkflowDto) (workflow *model.ProjectManagementWorkflow, err error) {
	if err = validateProjectManagementWorkflowDetails(input); err != nil {
		return nil, err
	}
	workflow = &model.ProjectManagementWorkflow{OrganizationId: claim.OrganizationId, Name: input.Name, Description: input.Description}
	workflow.Id = id
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		_, err := l.repos.ProjectManagementWorkflowRepo.LockTemplateWithTx(ctx, tx, id, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow not found")
		}
		if err := l.repos.ProjectManagementWorkflowRepo.UpdateTemplateWithTx(ctx, tx, workflow); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim.OrganizationId, id)
}

func (l *ProjectManagementWorkflowLogic) Delete(ctx context.Context, claim utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflows, err := l.repos.ProjectManagementWorkflowRepo.LockTemplatesWithTx(ctx, tx, claim.OrganizationId)
		if err != nil {
			return err
		}
		var workflow *model.ProjectManagementWorkflow
		enabledCount := 0
		for _, item := range workflows {
			if item.Id == id {
				workflow = item
			}
			if item.Enabled {
				enabledCount++
			}
		}
		if workflow == nil {
			return errorx.New("workflow not found")
		}
		if workflow.IsDefault {
			return errorx.New("default workflow cannot be deleted")
		}
		if workflow.Enabled && enabledCount <= 1 {
			return errorx.New("at least one enabled workflow is required")
		}
		if err := l.repos.ProjectManagementWorkflowRepo.DeleteTemplateWithTx(ctx, tx, id); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowDelete})
	})
}

func validateProjectManagementWorkflowDetails(input dto.ProjectManagementWorkflowDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("workflow name is required")
	}
	return nil
}

func cloneProjectManagementWorkflowDefinition(source *model.ProjectManagementWorkflow) ([]*model.ProjectManagementWorkflowStage, []*model.ProjectManagementWorkflowTransition) {
	stages := make([]*model.ProjectManagementWorkflowStage, len(source.Stages))
	for i, stage := range source.Stages {
		clone := *stage
		stages[i] = &clone
	}
	transitions := make([]*model.ProjectManagementWorkflowTransition, len(source.Transitions))
	for i, transition := range source.Transitions {
		clone := *transition
		transitions[i] = &clone
	}
	return stages, transitions
}

func projectManagementTemplateModels(input dto.ProjectManagementWorkflowDto) (*model.ProjectManagementWorkflow, []*model.ProjectManagementWorkflowStage, []*model.ProjectManagementWorkflowTransition, error) {
	workflow, err := copierx.Trans[dto.ProjectManagementWorkflowDto, model.ProjectManagementWorkflow](&input)
	if err != nil {
		return nil, nil, nil, err
	}
	stages, err := copierx.Trans[[]dto.ProjectManagementWorkflowStageDto, []*model.ProjectManagementWorkflowStage](&input.Stages)
	if err != nil {
		return nil, nil, nil, err
	}
	transitions, err := copierx.Trans[[]dto.ProjectManagementWorkflowTransitionDto, []*model.ProjectManagementWorkflowTransition](&input.Transitions)
	if err != nil {
		return nil, nil, nil, err
	}
	return workflow, *stages, *transitions, nil
}

func (l *ProjectManagementWorkflowLogic) recordOperation(ctx context.Context, tx *gorm.DB, input projectManagementWorkflowOperation) error {
	metadata, err := json.Marshal(map[string]int64{"workflow_id": input.workflowId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{
		EntityId:       input.claim.OrganizationId,
		EntityType:     constants.OperationEntityTypeProjectManagementOrganization,
		OperatorId:     input.claim.OrganizationUserId,
		OperatorName:   input.claim.Username,
		OperatorType:   constants.UserTypeBorrower,
		Operation:      input.action,
		Metadata:       metadata,
		OrganizationId: input.claim.OrganizationId,
	})
}

func isProjectManagementStagePhase(phase constants.ProjectManagementStagePhase) bool {
	return phase == constants.ProjectManagementStagePhasePending ||
		phase == constants.ProjectManagementStagePhaseInProgress ||
		phase == constants.ProjectManagementStagePhaseFinished
}

func isProjectManagementMemberRole(role constants.ProjectManagementMemberRole) bool {
	return role == constants.ProjectManagementMemberRoleProjectOwner ||
		role == constants.ProjectManagementMemberRoleProjectAdmin ||
		role == constants.ProjectManagementMemberRoleTaskCreator ||
		role == constants.ProjectManagementMemberRoleTaskExecutor
}
