package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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
	if err := validateProjectManagementWorkflowDefinition(workflow); err != nil {
		return err
	}
	return nil
}

func projectManagementWorkflowUpdateDefinition(input dto.ProjectManagementWorkflowUpdateDto) (dto.ProjectManagementWorkflowDto, error) {
	definition, err := copierx.Trans[dto.ProjectManagementWorkflowUpdateDto, dto.ProjectManagementWorkflowDto](&input)
	if err != nil {
		return dto.ProjectManagementWorkflowDto{}, err
	}
	stageIds := make(map[int64]struct{}, len(definition.Stages))
	for _, stage := range definition.Stages {
		if stage.Id != 0 {
			stageIds[stage.Id] = struct{}{}
		}
	}
	nextStageId := int64(-1)
	for i := range definition.Stages {
		if definition.Stages[i].Id != 0 {
			continue
		}
		for {
			if _, exists := stageIds[nextStageId]; !exists {
				definition.Stages[i].Id = nextStageId
				stageIds[nextStageId] = struct{}{}
				nextStageId--
				break
			}
			nextStageId--
		}
	}
	for i, stage := range input.Stages {
		if stage.Id == 0 && len(stage.Transitions) != 0 {
			return dto.ProjectManagementWorkflowDto{}, errorx.New("new workflow stage cannot include transitions")
		}
		for _, inputTransition := range stage.Transitions {
			if inputTransition.FromStageId != stage.Id {
				return dto.ProjectManagementWorkflowDto{}, errorx.New("workflow transition source stage is invalid")
			}
			transition, err := copierx.Trans[dto.ProjectManagementWorkflowUpdateTransitionDto, dto.ProjectManagementWorkflowTransitionDto](&inputTransition)
			if err != nil {
				return dto.ProjectManagementWorkflowDto{}, err
			}
			transition.FromStageId = definition.Stages[i].Id
			definition.Transitions = append(definition.Transitions, *transition)
		}
	}
	return *definition, nil
}

func validateProjectManagementWorkflowDefinition(workflow dto.ProjectManagementWorkflowDto) error {
	stageIDs := make(map[int64]struct{}, len(workflow.Stages))
	sortOrders := make(map[int]struct{}, len(workflow.Stages))
	pendingCount, finishedCount := 0, 0
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
		if stage.Phase == constants.ProjectManagementStagePhaseFinished {
			finishedCount++
		}
	}
	if pendingCount != 1 {
		return errorx.New("workflow must have exactly one pending stage")
	}
	if finishedCount == 0 {
		return errorx.New("workflow must have at least one finished stage")
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
	if _, err = l.List(ctx, claim.OrganizationId); err != nil {
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

func (l *ProjectManagementWorkflowLogic) Update(ctx context.Context, claim utils.JwtUser, id int64, input dto.ProjectManagementWorkflowUpdateDto) (workflow *model.ProjectManagementWorkflow, err error) {
	definition, err := projectManagementWorkflowUpdateDefinition(input)
	if err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		workflow, innerErr = l.repos.ProjectManagementWorkflowRepo.LockTemplateWithTx(ctx, tx, id, claim.OrganizationId)
		if innerErr != nil {
			return projectManagementNotFound(innerErr, "workflow not found")
		}
		if innerErr = l.assignNewStageSortOrders(ctx, tx, workflow.Id, &definition); innerErr != nil {
			return innerErr
		}
		if innerErr = ValidateProjectManagementWorkflow(definition); innerErr != nil {
			return innerErr
		}
		_, stages, transitions, innerErr := projectManagementTemplateModels(definition)
		if innerErr != nil {
			return innerErr
		}
		workflow.Name, workflow.Description = input.Name, input.Description
		workflow.RemoveRelationships()
		if _, innerErr = l.repos.ProjectManagementWorkflowRepo.UpdateOneByIdWithTx(ctx, tx, workflow.Id, workflow); innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementWorkflowRepo.ReplaceTemplateDefinitionWithTx(ctx, tx, workflow.Id, stages, transitions); innerErr != nil {
			return innerErr
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim.OrganizationId, id)
}

func (l *ProjectManagementWorkflowLogic) assignNewStageSortOrders(ctx context.Context, tx *gorm.DB, workflowId int64, definition *dto.ProjectManagementWorkflowDto) error {
	maxSortOrder, err := l.repos.ProjectManagementWorkflowRepo.MaxStageSortOrderWithTx(ctx, tx, workflowId)
	if err != nil {
		return err
	}
	for i := range definition.Stages {
		if definition.Stages[i].Id > 0 {
			continue
		}
		maxSortOrder++
		definition.Stages[i].SortOrder = maxSortOrder
	}
	return nil
}

func (l *ProjectManagementWorkflowLogic) SetDefault(ctx context.Context, claim utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflows, err := l.repos.ProjectManagementWorkflowRepo.LockTemplatesWithTx(ctx, tx, claim.OrganizationId)
		if err != nil {
			return err
		}
		for _, workflow := range workflows {
			if workflow.Id != id {
				continue
			}
			if workflow.IsDefault {
				return nil
			}
			if err = l.repos.ProjectManagementWorkflowRepo.DemoteDefaultWithTx(ctx, tx, claim.OrganizationId, id); err != nil {
				return err
			}
			if err = l.repos.ProjectManagementWorkflowRepo.SetDefaultWithTx(ctx, tx, id, claim.OrganizationId); err != nil {
				return err
			}
			return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowUpdate})
		}
		return errorx.New("workflow not found")
	})
}

func (l *ProjectManagementWorkflowLogic) Delete(ctx context.Context, claim utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflows, err := l.repos.ProjectManagementWorkflowRepo.LockTemplatesWithTx(ctx, tx, claim.OrganizationId)
		if err != nil {
			return err
		}
		var workflow *model.ProjectManagementWorkflow
		for _, item := range workflows {
			if item.Id == id {
				workflow = item
			}
		}
		if workflow == nil {
			return errorx.New("workflow not found")
		}
		if workflow.IsDefault {
			return errorx.New("default workflow cannot be deleted")
		}
		if err := l.repos.ProjectManagementWorkflowRepo.DeleteTemplateWithTx(ctx, tx, id); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowDelete})
	})
}

func (l *ProjectManagementWorkflowLogic) StageAll(ctx context.Context, organizationId, workflowId int64) ([]*model.ProjectManagementWorkflowStage, error) {
	if _, err := l.Detail(ctx, organizationId, workflowId); err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementWorkflowRepo.ListTemplateStages(ctx, workflowId, organizationId)
}

func (l *ProjectManagementWorkflowLogic) StageDetail(ctx context.Context, organizationId, workflowId, id int64) (*model.ProjectManagementWorkflowStage, error) {
	stage, err := l.repos.ProjectManagementWorkflowRepo.FindTemplateStage(ctx, id, workflowId, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "workflow stage not found")
	}
	return stage, nil
}

func (l *ProjectManagementWorkflowLogic) StageCreate(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementWorkflowStageDto, workflowId int64) (stage *model.ProjectManagementWorkflowStage, err error) {
	if err = validateProjectManagementWorkflowStage(input); err != nil {
		return nil, err
	}
	stage, err = copierx.Trans[dto.ProjectManagementWorkflowStageDto, model.ProjectManagementWorkflowStage](&input)
	if err != nil {
		return nil, err
	}
	stage.WorkflowId = workflowId
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err = l.repos.ProjectManagementWorkflowRepo.LockTemplateWithTx(ctx, tx, workflowId, claim.OrganizationId); err != nil {
			return projectManagementNotFound(err, "workflow not found")
		}
		maxSortOrder, err := l.repos.ProjectManagementWorkflowRepo.MaxStageSortOrderWithTx(ctx, tx, workflowId)
		if err != nil {
			return err
		}
		stage.SortOrder = maxSortOrder + 1
		if err = l.repos.ProjectManagementWorkflowRepo.CreateTemplateStageWithTx(ctx, tx, stage); err != nil {
			return err
		}
		if err = l.validateWorkflowStages(ctx, tx, workflowId); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflowId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
	return stage, err
}

func (l *ProjectManagementWorkflowLogic) StageSwapOrder(ctx context.Context, claim utils.JwtUser, workflowId, id, targetId int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if innerErr := l.repos.ProjectManagementWorkflowRepo.SwapTemplateStageSortOrdersWithTx(ctx, tx, workflowId, claim.OrganizationId, id, targetId); innerErr != nil {
			return innerErr
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflowId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
}

func (l *ProjectManagementWorkflowLogic) StageUpdate(ctx context.Context, claim utils.JwtUser, workflowId, id int64, input dto.ProjectManagementWorkflowStageDto) (stage *model.ProjectManagementWorkflowStage, err error) {
	if err = validateProjectManagementWorkflowStage(input); err != nil {
		return nil, err
	}
	stage, err = copierx.Trans[dto.ProjectManagementWorkflowStageDto, model.ProjectManagementWorkflowStage](&input)
	if err != nil {
		return nil, err
	}
	stage.Id = id
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		existing, err := l.repos.ProjectManagementWorkflowRepo.FindTemplateStageWithTx(ctx, tx, id, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow stage not found")
		}
		stage.WorkflowId = existing.WorkflowId
		if err = l.repos.ProjectManagementWorkflowRepo.UpdateTemplateStageWithTx(ctx, tx, stage); err != nil {
			return err
		}
		if err = l.validateWorkflowStages(ctx, tx, stage.WorkflowId); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: stage.WorkflowId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
	if err != nil {
		return nil, err
	}
	return l.StageDetail(ctx, claim.OrganizationId, workflowId, id)
}

func (l *ProjectManagementWorkflowLogic) StageDelete(ctx context.Context, claim utils.JwtUser, workflowId, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindTemplateStageWithTx(ctx, tx, id, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow stage not found")
		}
		occupied, err := l.repos.ProjectManagementTaskRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]any{"stage_id": stage.Id})
		if err != nil {
			return err
		}
		if occupied {
			return errorx.New("workflow stage is referenced by a task")
		}
		if err := l.repos.ProjectManagementWorkflowRepo.DeleteTemplateStageWithTx(ctx, tx, stage); err != nil {
			return err
		}
		if err := l.validateWorkflowStages(ctx, tx, stage.WorkflowId); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: stage.WorkflowId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
}

func (l *ProjectManagementWorkflowLogic) TransitionAll(ctx context.Context, organizationId, workflowId, fromStageId int64) ([]*model.ProjectManagementWorkflowTransition, error) {
	if _, err := l.StageDetail(ctx, organizationId, workflowId, fromStageId); err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementWorkflowRepo.ListTransitionsByFromStage(ctx, fromStageId)
}

func (l *ProjectManagementWorkflowLogic) TransitionDetail(ctx context.Context, organizationId, workflowId, id int64) (*model.ProjectManagementWorkflowTransition, error) {
	transition, err := l.repos.ProjectManagementWorkflowRepo.FindTemplateTransition(ctx, id, workflowId, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "workflow transition not found")
	}
	return transition, nil
}

func (l *ProjectManagementWorkflowLogic) TransitionCreate(ctx context.Context, claim utils.JwtUser, workflowId int64, input dto.ProjectManagementWorkflowTransitionDto) (transition *model.ProjectManagementWorkflowTransition, err error) {
	if err = validateProjectManagementWorkflowTransition(input); err != nil {
		return nil, err
	}
	transition, err = copierx.Trans[dto.ProjectManagementWorkflowTransitionDto, model.ProjectManagementWorkflowTransition](&input)
	if err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		fromStage, err := l.repos.ProjectManagementWorkflowRepo.FindTemplateStageWithTx(ctx, tx, input.FromStageId, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow stage not found")
		}
		toStage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, input.ToStageId)
		if err != nil || toStage.WorkflowId != fromStage.WorkflowId {
			return errorx.New("workflow transition target stage is invalid")
		}
		if _, err = l.repos.ProjectManagementWorkflowRepo.FindTransitionWithTx(ctx, tx, fromStage.WorkflowId, input.FromStageId, input.ToStageId); err == nil {
			return errorx.New("workflow transition already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		transition.WorkflowId = fromStage.WorkflowId
		if err = l.repos.ProjectManagementWorkflowRepo.CreateTemplateTransitionWithTx(ctx, tx, transition); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: transition.WorkflowId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
	return transition, err
}

func (l *ProjectManagementWorkflowLogic) TransitionUpdate(ctx context.Context, claim utils.JwtUser, workflowId, id int64, input dto.ProjectManagementWorkflowTransitionDto) (transition *model.ProjectManagementWorkflowTransition, err error) {
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		transition, err = l.repos.ProjectManagementWorkflowRepo.FindTemplateTransitionWithTx(ctx, tx, id, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow transition not found")
		}
		if err = validateProjectManagementWorkflowTransition(input); err != nil {
			return err
		}
		fromStage, err := l.repos.ProjectManagementWorkflowRepo.FindTemplateStageWithTx(ctx, tx, input.FromStageId, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow stage not found")
		}
		toStage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, input.ToStageId)
		if err != nil || toStage.WorkflowId != fromStage.WorkflowId {
			return errorx.New("workflow transition target stage is invalid")
		}
		duplicate, duplicateErr := l.repos.ProjectManagementWorkflowRepo.FindTransitionWithTx(ctx, tx, transition.WorkflowId, input.FromStageId, input.ToStageId)
		if duplicateErr == nil && duplicate.Id != id {
			return errorx.New("workflow transition already exists")
		}
		if duplicateErr != nil && !errors.Is(duplicateErr, gorm.ErrRecordNotFound) {
			return duplicateErr
		}
		transitionId := transition.Id
		if err = copierx.Assign(&input, transition); err != nil {
			return err
		}
		transition.Id = transitionId
		if err = l.repos.ProjectManagementWorkflowRepo.UpdateTemplateTransitionWithTx(ctx, tx, transition); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: transition.WorkflowId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
	if err != nil {
		return nil, err
	}
	return l.TransitionDetail(ctx, claim.OrganizationId, workflowId, id)
}

func (l *ProjectManagementWorkflowLogic) TransitionDelete(ctx context.Context, claim utils.JwtUser, workflowId, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		transition, err := l.repos.ProjectManagementWorkflowRepo.FindTemplateTransitionWithTx(ctx, tx, id, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow transition not found")
		}
		if err := l.repos.ProjectManagementWorkflowRepo.DeleteTemplateTransitionWithTx(ctx, tx, transition); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: transition.WorkflowId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
}

func (l *ProjectManagementWorkflowLogic) validateWorkflowStages(ctx context.Context, tx *gorm.DB, workflowId int64) error {
	stages, err := l.repos.ProjectManagementWorkflowRepo.ListStagesWithTx(ctx, tx, workflowId)
	if err != nil {
		return err
	}
	pendingCount, finishedCount := 0, 0
	for _, stage := range stages {
		if !isProjectManagementStagePhase(stage.Phase) {
			return errorx.New("invalid workflow stage phase")
		}
		switch stage.Phase {
		case constants.ProjectManagementStagePhasePending:
			pendingCount++
		case constants.ProjectManagementStagePhaseFinished:
			finishedCount++
		}
	}
	if pendingCount != 1 {
		return errorx.New("workflow must have exactly one pending stage")
	}
	if finishedCount == 0 {
		return errorx.New("workflow must have at least one finished stage")
	}
	return nil
}

func validateProjectManagementWorkflowStage(input dto.ProjectManagementWorkflowStageDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("workflow stage name is required")
	}
	if !isProjectManagementStagePhase(input.Phase) {
		return errorx.New("invalid workflow stage phase")
	}
	if input.SortOrder < 0 {
		return errorx.New("workflow stage sort order must be >= 0")
	}
	return nil
}

func validateProjectManagementWorkflowTransition(input dto.ProjectManagementWorkflowTransitionDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("workflow transition name is required")
	}
	if input.FromStageId == 0 || input.ToStageId == 0 || input.FromStageId == input.ToStageId {
		return errorx.New("workflow transition stages are invalid")
	}
	roles := make(map[constants.ProjectManagementMemberRole]struct{}, len(input.AllowedActorTypes))
	for _, role := range input.AllowedActorTypes {
		if !isProjectManagementMemberRole(role) {
			return errorx.New("invalid workflow transition claim type")
		}
		if _, exists := roles[role]; exists {
			return errorx.New("workflow transition claim type must be unique")
		}
		roles[role] = struct{}{}
	}
	return nil
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
