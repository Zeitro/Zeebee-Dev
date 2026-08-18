package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

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
		for _, actor := range transition.AllowedActorTypes {
			if !isProjectManagementMemberRole(actor) {
				return errorx.New("invalid workflow transition actor type")
			}
			if _, exists := actors[actor]; exists {
				return errorx.New("workflow transition actor type must be unique")
			}
			actors[actor] = struct{}{}
		}
	}
	return nil
}

func (l *ProjectManagementLogic) PageWorkflows(ctx context.Context, organizationId int64, input dto.ProjectManagementWorkflowPageRequestDto) (gormx.PageResponse[*model.ProjectManagementWorkflow], error) {
	return l.repos.ProjectManagementWorkflowRepo.Page(ctx, organizationId, gormx.PageRequest{Page: input.Page, Limit: input.Limit}, input.Filter.Keyword, input.Filter.Enabled)
}

func (l *ProjectManagementLogic) GetWorkflow(ctx context.Context, organizationId, id int64) (*model.ProjectManagementWorkflow, error) {
	workflow, err := l.repos.ProjectManagementWorkflowRepo.FindDefinitionByOrganization(ctx, id, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "workflow not found")
	}
	return workflow, nil
}

func (l *ProjectManagementLogic) CreateWorkflow(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementWorkflowDto) (workflow *model.ProjectManagementWorkflow, err error) {
	normalizeProjectManagementTemplateStageReferences(&input)
	if err = validateProjectManagementTemplateWorkflow(input); err != nil {
		return nil, err
	}
	workflow, stages, transitions, err := projectManagementTemplateModels(input)
	if err != nil {
		return nil, err
	}
	workflow.Kind = constants.ProjectManagementWorkflowKindOrganizationTemplate
	workflow.OrganizationId = actor.OrganizationId
	workflow.ProjectId = 0
	workflow.SourceTemplateId = 0
	workflow.SourceTemplateName = ""
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if workflow.IsDefault {
			if err := l.repos.ProjectManagementWorkflowRepo.DemoteDefaultWithTx(ctx, tx, actor.OrganizationId, 0); err != nil {
				return err
			}
		}
		workflow, err = l.repos.ProjectManagementWorkflowRepo.CreateWithTx(ctx, tx, workflow, false)
		if err != nil {
			return err
		}
		if err := l.repos.ProjectManagementWorkflowRepo.ReplaceTemplateDefinitionWithTx(ctx, tx, workflow.Id, stages, transitions); err != nil {
			return err
		}
		return l.recordWorkflowOperation(ctx, tx, actor, workflow.Id, constants.OperationActionProjectManagementWorkflowCreate)
	})
	if err != nil {
		return nil, err
	}
	return l.GetWorkflow(ctx, actor.OrganizationId, workflow.Id)
}

func (l *ProjectManagementLogic) UpdateWorkflow(ctx context.Context, actor *dto.ProjectManagementActor, id int64, input dto.ProjectManagementWorkflowDto) (workflow *model.ProjectManagementWorkflow, err error) {
	if err = validateProjectManagementTemplateWorkflow(input); err != nil {
		return nil, err
	}
	workflow, stages, transitions, err := projectManagementTemplateModels(input)
	if err != nil {
		return nil, err
	}
	workflow.Id = id
	workflow.OrganizationId = actor.OrganizationId
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		existing, err := l.repos.ProjectManagementWorkflowRepo.LockTemplateWithTx(ctx, tx, id, actor.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow not found")
		}
		if existing.IsDefault && !workflow.IsDefault {
			return errorx.New("default workflow must be replaced by another default workflow")
		}
		if workflow.IsDefault {
			if err := l.repos.ProjectManagementWorkflowRepo.DemoteDefaultWithTx(ctx, tx, actor.OrganizationId, id); err != nil {
				return err
			}
		}
		if err := l.repos.ProjectManagementWorkflowRepo.UpdateTemplateWithTx(ctx, tx, workflow); err != nil {
			return err
		}
		if err := l.repos.ProjectManagementWorkflowRepo.ReplaceTemplateDefinitionWithTx(ctx, tx, id, stages, transitions); err != nil {
			return err
		}
		return l.recordWorkflowOperation(ctx, tx, actor, id, constants.OperationActionProjectManagementWorkflowUpdate)
	})
	if err != nil {
		return nil, err
	}
	return l.GetWorkflow(ctx, actor.OrganizationId, id)
}

func (l *ProjectManagementLogic) DeleteWorkflow(ctx context.Context, actor *dto.ProjectManagementActor, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflow, err := l.repos.ProjectManagementWorkflowRepo.LockTemplateWithTx(ctx, tx, id, actor.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow not found")
		}
		if workflow.IsDefault {
			return errorx.New("default workflow cannot be deleted")
		}
		if err := l.repos.ProjectManagementWorkflowRepo.DeleteTemplateWithTx(ctx, tx, id); err != nil {
			return err
		}
		return l.recordWorkflowOperation(ctx, tx, actor, id, constants.OperationActionProjectManagementWorkflowDelete)
	})
}

func validateProjectManagementTemplateWorkflow(input dto.ProjectManagementWorkflowDto) error {
	if input.IsDefault && !input.Enabled {
		return errorx.New("default workflow cannot be disabled")
	}
	return ValidateProjectManagementWorkflow(input)
}

func normalizeProjectManagementTemplateStageReferences(input *dto.ProjectManagementWorkflowDto) {
	used := make(map[int64]struct{}, len(input.Stages))
	for _, stage := range input.Stages {
		if stage.Id != 0 {
			used[stage.Id] = struct{}{}
		}
	}

	referenceId := int64(-1)
	for i := range input.Stages {
		if input.Stages[i].Id != 0 {
			continue
		}
		for {
			if _, exists := used[referenceId]; !exists {
				break
			}
			referenceId--
		}
		input.Stages[i].Id = referenceId
		used[referenceId] = struct{}{}
		referenceId--
	}
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
	for i, transition := range input.Transitions {
		(*transitions)[i].AllowedActorTypes = pq.StringArray(make([]string, len(transition.AllowedActorTypes)))
		for j, actorType := range transition.AllowedActorTypes {
			(*transitions)[i].AllowedActorTypes[j] = string(actorType)
		}
	}
	return workflow, *stages, *transitions, nil
}

func (l *ProjectManagementLogic) recordWorkflowOperation(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, workflowId int64, action constants.OperationAction) error {
	metadata, err := json.Marshal(map[string]int64{"workflow_id": workflowId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{
		EntityId:       actor.OrganizationId,
		EntityType:     constants.OperationEntityTypeProjectManagementOrganization,
		OperatorId:     actor.OrganizationUserId,
		OperatorName:   actor.Name,
		OperatorType:   constants.UserTypeBorrower,
		Operation:      action,
		Metadata:       metadata,
		OrganizationId: actor.OrganizationId,
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
