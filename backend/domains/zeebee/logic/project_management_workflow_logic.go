package logic

import (
	"context"
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
	claim         utils.JwtUser
	workflowId    int64
	stageId       int64
	targetStageId int64
	transitionId  int64
	action        constants.OperationAction
	before        map[string]any
	after         map[string]any
}

func NewProjectManagementWorkflowLogic(baseLogic *BaseLogic) *ProjectManagementWorkflowLogic {
	return &ProjectManagementWorkflowLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementWorkflowLogic) List(ctx context.Context, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows, err := l.repos.ProjectManagementWorkflowRepo.ListByOrganization(ctx, organizationId)
	if err != nil {
		return workflows, err
	}
	for _, workflow := range workflows {
		if workflow.IsDefault {
			return workflows, nil
		}
	}
	if err = l.ensureDefault(ctx, organizationId); err != nil {
		workflows, retryErr := l.repos.ProjectManagementWorkflowRepo.ListByOrganization(ctx, organizationId)
		if retryErr == nil {
			for _, workflow := range workflows {
				if workflow.IsDefault {
					return workflows, nil
				}
			}
		}
		return nil, err
	}
	return l.repos.ProjectManagementWorkflowRepo.ListByOrganization(ctx, organizationId)
}

func (l *ProjectManagementWorkflowLogic) ensureDefault(ctx context.Context, organizationId int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflow, err := l.repos.ProjectManagementWorkflowRepo.CreateWithTx(ctx, tx, &model.ProjectManagementWorkflow{
			Kind:           constants.ProjectManagementWorkflowKindOrganizationTemplate,
			OrganizationId: organizationId,
			Name:           "Default",
			IsDefault:      true,
		}, false)
		if err != nil {
			return err
		}
		stages := []*model.ProjectManagementWorkflowStage{
			{WorkflowId: workflow.Id, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, Color: "#3b82f6", SortOrder: 10},
			{WorkflowId: workflow.Id, Name: "In progress", Phase: constants.ProjectManagementStagePhaseInProgress, Color: "#f59e0b", SortOrder: 20},
			{WorkflowId: workflow.Id, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, Color: "#22c55e", SortOrder: 30},
		}
		if err = tx.WithContext(ctx).Create(&stages).Error; err != nil {
			return errorx.WithStack(err)
		}
		allowedActorTypes := []constants.ProjectManagementMemberRole{
			constants.ProjectManagementMemberRoleProjectOwner,
			constants.ProjectManagementMemberRoleProjectAdmin,
			constants.ProjectManagementMemberRoleTaskExecutor,
			constants.ProjectManagementMemberRoleTaskCreator,
		}
		transitions := make([]*model.ProjectManagementWorkflowTransition, 0, len(stages)*(len(stages)-1))
		for _, from := range stages {
			for _, to := range stages {
				if from.Id != to.Id {
					transitions = append(transitions, &model.ProjectManagementWorkflowTransition{WorkflowId: workflow.Id, FromStageId: from.Id, ToStageId: to.Id, AllowedActorTypes: allowedActorTypes})
				}
			}
		}
		return errorx.WithStack(tx.WithContext(ctx).Create(&transitions).Error)
	})
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
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflow.Id, action: constants.OperationActionProjectManagementWorkflowCreate, after: map[string]any{"id": workflow.Id, "name": workflow.Name}})
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
		before := workflowOperationValues(workflow)
		if innerErr = l.assignNewStageSortOrders(ctx, tx, workflow.Id, &definition); innerErr != nil {
			return innerErr
		}
		if innerErr = ValidateProjectManagementWorkflow(definition); innerErr != nil {
			return innerErr
		}
		beforeStages, innerErr := l.repos.ProjectManagementWorkflowRepo.ListStagesWithTx(ctx, tx, workflow.Id)
		if innerErr != nil {
			return innerErr
		}
		beforeTransitions, innerErr := l.repos.ProjectManagementWorkflowRepo.ListTransitionsWithTx(ctx, tx, workflow.Id)
		if innerErr != nil {
			return innerErr
		}
		_, stages, transitions, innerErr := projectManagementTemplateModels(definition)
		if innerErr != nil {
			return innerErr
		}
		desiredStages, desiredTransitions := cloneProjectManagementWorkflowDefinition(&model.ProjectManagementWorkflow{Stages: stages, Transitions: transitions})
		workflow.Name, workflow.Description = input.Name, input.Description
		workflow.RemoveRelationships()
		if _, innerErr = l.repos.ProjectManagementWorkflowRepo.UpdateOneByIdWithTx(ctx, tx, workflow.Id, workflow); innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementWorkflowRepo.ReplaceTemplateDefinitionWithTx(ctx, tx, workflow.Id, stages, transitions); innerErr != nil {
			return innerErr
		}
		after := workflowOperationValues(workflow)
		before, after = operationValuesDiff(before, after)
		if len(before) != 0 {
			if innerErr = l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowUpdate, before: before, after: after}); innerErr != nil {
				return innerErr
			}
		}
		return l.recordTemplateWorkflowDefinitionChanges(ctx, tx, claim, id, beforeStages, beforeTransitions, desiredStages, desiredTransitions, stages, transitions)
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim.OrganizationId, id)
}

func (l *ProjectManagementWorkflowLogic) recordTemplateWorkflowDefinitionChanges(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, workflowId int64, beforeStages []*model.ProjectManagementWorkflowStage, beforeTransitions []*model.ProjectManagementWorkflowTransition, desiredStages []*model.ProjectManagementWorkflowStage, desiredTransitions []*model.ProjectManagementWorkflowTransition, stages []*model.ProjectManagementWorkflowStage, transitions []*model.ProjectManagementWorkflowTransition) error {
	record := func(action constants.OperationAction, stageId, transitionId int64, before, after map[string]any) error {
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflowId, stageId: stageId, transitionId: transitionId, action: action, before: before, after: after})
	}
	beforeStageByID := make(map[int64]*model.ProjectManagementWorkflowStage, len(beforeStages))
	for _, stage := range beforeStages {
		beforeStageByID[stage.Id] = stage
	}
	for i, desired := range desiredStages {
		if before := beforeStageByID[desired.Id]; before != nil {
			old, next := operationValuesDiff(workflowStageOperationValues(before), workflowStageOperationValues(desired))
			if len(old) != 0 {
				if err := record(constants.OperationActionProjectManagementWorkflowStageUpdate, desired.Id, 0, old, next); err != nil {
					return err
				}
			}
			delete(beforeStageByID, desired.Id)
			continue
		}
		if err := record(constants.OperationActionProjectManagementWorkflowStageCreate, stages[i].Id, 0, nil, workflowStageOperationIdentity(stages[i])); err != nil {
			return err
		}
	}
	for _, stage := range beforeStageByID {
		if err := record(constants.OperationActionProjectManagementWorkflowStageDelete, stage.Id, 0, workflowStageOperationIdentity(stage), nil); err != nil {
			return err
		}
	}
	beforeTransitionByID := make(map[int64]*model.ProjectManagementWorkflowTransition, len(beforeTransitions))
	for _, transition := range beforeTransitions {
		beforeTransitionByID[transition.Id] = transition
	}
	for i, desired := range desiredTransitions {
		if before := beforeTransitionByID[desired.Id]; before != nil {
			old, next := operationValuesDiff(workflowTransitionOperationValues(before), workflowTransitionOperationValues(desired))
			if len(old) != 0 {
				if err := record(constants.OperationActionProjectManagementWorkflowTransitionUpdate, 0, desired.Id, old, next); err != nil {
					return err
				}
			}
			delete(beforeTransitionByID, desired.Id)
			continue
		}
		if err := record(constants.OperationActionProjectManagementWorkflowTransitionCreate, 0, transitions[i].Id, nil, workflowTransitionOperationIdentity(transitions[i])); err != nil {
			return err
		}
	}
	for _, transition := range beforeTransitionByID {
		if err := record(constants.OperationActionProjectManagementWorkflowTransitionDelete, 0, transition.Id, workflowTransitionOperationIdentity(transition), nil); err != nil {
			return err
		}
	}
	return nil
}

func (l *ProjectManagementWorkflowLogic) assignNewStageSortOrders(ctx context.Context, tx *gorm.DB, workflowId int64, definition *dto.ProjectManagementWorkflowDto) error {
	maxSortOrder, err := l.repos.ProjectManagementWorkflowStageRepo.MaxSortOrderWithTx(ctx, tx, workflowId)
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
			return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowUpdate, before: map[string]any{"is_default": false}, after: map[string]any{"is_default": true}})
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
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowDelete, before: workflowOperationValues(workflow)})
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
		maxSortOrder, err := l.repos.ProjectManagementWorkflowStageRepo.MaxSortOrderWithTx(ctx, tx, workflowId)
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
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflowId, stageId: stage.Id, action: constants.OperationActionProjectManagementWorkflowStageCreate, after: workflowStageOperationIdentity(stage)})
	})
	return stage, err
}

func (l *ProjectManagementWorkflowLogic) StageSwapOrder(ctx context.Context, claim utils.JwtUser, workflowId, id, targetId int64) error {
	if id == targetId {
		return nil
	}
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindSwappableStageWithTx(ctx, tx, id, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow stage not found")
		}
		target, err := l.repos.ProjectManagementWorkflowRepo.FindSwappableStageWithTx(ctx, tx, targetId, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow stage not found")
		}
		if innerErr := l.repos.ProjectManagementWorkflowRepo.SwapStageSortOrdersWithTx(ctx, tx, workflowId, claim.OrganizationId, id, targetId); innerErr != nil {
			return innerErr
		}
		before := map[string]any{"sort_order": stage.SortOrder, "target_sort_order": target.SortOrder}
		after := map[string]any{"sort_order": target.SortOrder, "target_sort_order": stage.SortOrder}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflowId, stageId: stage.Id, targetStageId: target.Id, action: constants.OperationActionProjectManagementWorkflowStageUpdate, before: before, after: after})
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
		before := workflowStageOperationValues(existing)
		stage.WorkflowId = existing.WorkflowId
		if err = l.repos.ProjectManagementWorkflowRepo.UpdateTemplateStageWithTx(ctx, tx, stage); err != nil {
			return err
		}
		if err = l.validateWorkflowStages(ctx, tx, stage.WorkflowId); err != nil {
			return err
		}
		after := workflowStageOperationValues(stage)
		before, after = operationValuesDiff(before, after)
		if len(before) == 0 {
			return nil
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: stage.WorkflowId, stageId: stage.Id, action: constants.OperationActionProjectManagementWorkflowStageUpdate, before: before, after: after})
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
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: stage.WorkflowId, stageId: stage.Id, action: constants.OperationActionProjectManagementWorkflowStageDelete, before: workflowStageOperationIdentity(stage)})
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
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: transition.WorkflowId, transitionId: transition.Id, action: constants.OperationActionProjectManagementWorkflowTransitionCreate, after: workflowTransitionOperationIdentity(transition)})
	})
	return transition, err
}

func (l *ProjectManagementWorkflowLogic) TransitionUpdate(ctx context.Context, claim utils.JwtUser, workflowId, id int64, input dto.ProjectManagementWorkflowTransitionDto) (transition *model.ProjectManagementWorkflowTransition, err error) {
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		transition, err = l.repos.ProjectManagementWorkflowRepo.FindTemplateTransitionWithTx(ctx, tx, id, workflowId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow transition not found")
		}
		before := workflowTransitionOperationValues(transition)
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
		after := workflowTransitionOperationValues(transition)
		before, after = operationValuesDiff(before, after)
		if len(before) == 0 {
			return nil
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: transition.WorkflowId, transitionId: transition.Id, action: constants.OperationActionProjectManagementWorkflowTransitionUpdate, before: before, after: after})
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
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: transition.WorkflowId, transitionId: transition.Id, action: constants.OperationActionProjectManagementWorkflowTransitionDelete, before: workflowTransitionOperationIdentity(transition)})
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
	if input.FromStageId == 0 || input.ToStageId == 0 {
		return errorx.New("workflow transition source and target stage IDs are required")
	}
	if input.FromStageId == input.ToStageId {
		return errorx.New(fmt.Sprintf("workflow transition source and target stages must differ (both are %d)", input.FromStageId))
	}
	roles := make(map[constants.ProjectManagementMemberRole]struct{}, len(input.AllowedActorTypes))
	for _, role := range input.AllowedActorTypes {
		if !isProjectManagementMemberRole(role) {
			return errorx.New(fmt.Sprintf("workflow transition allowed actor type %q is invalid", role))
		}
		if _, exists := roles[role]; exists {
			return errorx.New(fmt.Sprintf("workflow transition allowed actor type %q is duplicated", role))
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
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{
		WorkflowID:    input.workflowId,
		StageID:       input.stageId,
		TargetStageID: input.targetStageId,
		TransitionID:  input.transitionId,
	}, input.before, input.after)
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

func workflowOperationValues(workflow *model.ProjectManagementWorkflow) map[string]any {
	return map[string]any{"name": workflow.Name, "description": workflow.Description, "is_default": workflow.IsDefault}
}

func workflowStageOperationValues(stage *model.ProjectManagementWorkflowStage) map[string]any {
	return map[string]any{"name": stage.Name, "description": stage.Description, "phase": stage.Phase, "color": stage.Color, "sort_order": stage.SortOrder}
}

func workflowStageOperationIdentity(stage *model.ProjectManagementWorkflowStage) map[string]any {
	return map[string]any{"id": stage.Id, "name": stage.Name}
}

func workflowTransitionOperationValues(transition *model.ProjectManagementWorkflowTransition) map[string]any {
	return map[string]any{"name": transition.Name, "from_stage_id": transition.FromStageId, "to_stage_id": transition.ToStageId, "allowed_actor_types": transition.AllowedActorTypes}
}

func workflowTransitionOperationIdentity(transition *model.ProjectManagementWorkflowTransition) map[string]any {
	return map[string]any{"id": transition.Id, "name": transition.Name}
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
