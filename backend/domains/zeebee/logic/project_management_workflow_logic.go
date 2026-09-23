package logic

import (
	"context"
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
	before     map[string]any
	after      map[string]any
}

func NewProjectManagementWorkflowLogic(baseLogic *BaseLogic) *ProjectManagementWorkflowLogic {
	return &ProjectManagementWorkflowLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementWorkflowLogic) List(ctx context.Context, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows, err := l.repos.ProjectManagementWorkflowRepo.ListWorkflowSources(ctx, organizationId)
	if err != nil || hasDefaultWorkflow(workflows) {
		return workflows, err
	}
	if err = l.ensureDefault(ctx, organizationId); err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementWorkflowRepo.ListWorkflowSources(ctx, organizationId)
}

func hasDefaultWorkflow(workflows []*model.ProjectManagementWorkflow) bool {
	for _, workflow := range workflows {
		if workflow.IsDefault {
			return true
		}
	}
	return false
}

func (l *ProjectManagementWorkflowLogic) ensureDefault(ctx context.Context, organizationId int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflows, err := l.repos.ProjectManagementWorkflowRepo.LockWorkflowSourcesWithTx(ctx, tx, organizationId)
		if err != nil || hasDefaultWorkflow(workflows) {
			return err
		}
		template := &model.ProjectManagementTemplate{TemplateType: constants.ProjectManagementTemplateTypeOrganization, OrganizationId: organizationId, Name: "Default", Status: "Active", IsDefault: true}
		if err = tx.WithContext(ctx).Create(template).Error; err != nil {
			return errorx.WithStack(err)
		}
		workflow := &model.ProjectManagementWorkflow{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: template.Id, Name: template.Name, IsDefault: true}
		if err = tx.WithContext(ctx).Omit("Kind", "OrganizationId", "ProjectId", "SourceTemplateId", "SourceTemplateName").Create(workflow).Error; err != nil {
			return errorx.WithStack(err)
		}
		return l.replaceWorkflowSourceDefinition(ctx, tx, workflow, defaultWorkflowDefinition())
	})
}

func defaultWorkflowDefinition() dto.ProjectManagementWorkflowDto {
	return dto.ProjectManagementWorkflowDto{
		Stages: []dto.ProjectManagementWorkflowStageDto{
			{Id: -1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, Color: "#3b82f6", SortOrder: 10},
			{Id: -2, Name: "In progress", Phase: constants.ProjectManagementStagePhaseInProgress, Color: "#f59e0b", SortOrder: 20},
			{Id: -3, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, Color: "#22c55e", SortOrder: 30},
		},
		Transitions: []dto.ProjectManagementWorkflowTransitionDto{
			{Id: -1, Name: "Start", FromStageId: -1, ToStageId: -2, AllowedActorTypes: defaultWorkflowActorTypes()},
			{Id: -2, Name: "Finish", FromStageId: -2, ToStageId: -3, AllowedActorTypes: defaultWorkflowActorTypes()},
		},
	}
}

func defaultWorkflowActorTypes() []constants.ProjectManagementMemberRole {
	return []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin, constants.ProjectManagementMemberRoleTaskCreator, constants.ProjectManagementMemberRoleTaskExecutor}
}

func ValidateProjectManagementWorkflow(workflow dto.ProjectManagementWorkflowDto) error {
	if strings.TrimSpace(workflow.Name) == "" {
		return errorx.New("workflow name is required")
	}
	return validateProjectManagementWorkflowDefinition(workflow)
}

func projectManagementWorkflowUpdateDefinition(input dto.ProjectManagementWorkflowUpdateDto) (dto.ProjectManagementWorkflowDto, error) {
	definition, err := copierx.Trans[dto.ProjectManagementWorkflowUpdateDto, dto.ProjectManagementWorkflowDto](&input)
	if err != nil {
		return dto.ProjectManagementWorkflowDto{}, err
	}
	return normalizeWorkflowDefinition(*definition, input.Stages)
}

func normalizeWorkflowDefinition(definition dto.ProjectManagementWorkflowDto, input []dto.ProjectManagementWorkflowUpdateStageDto) (dto.ProjectManagementWorkflowDto, error) {
	stageIDs := make(map[int64]struct{}, len(definition.Stages))
	for _, stage := range definition.Stages {
		if stage.Id != 0 {
			stageIDs[stage.Id] = struct{}{}
		}
	}
	nextID := int64(-1)
	for i := range definition.Stages {
		if definition.Stages[i].Id != 0 {
			continue
		}
		for {
			if _, exists := stageIDs[nextID]; !exists {
				definition.Stages[i].Id = nextID
				stageIDs[nextID] = struct{}{}
				nextID--
				break
			}
			nextID--
		}
	}
	for i, stage := range input {
		if stage.Id == 0 && len(stage.Transitions) != 0 {
			return dto.ProjectManagementWorkflowDto{}, errorx.New("new workflow stage cannot include transitions")
		}
		for _, item := range stage.Transitions {
			if item.FromStageId != stage.Id {
				return dto.ProjectManagementWorkflowDto{}, errorx.New("workflow transition source stage is invalid")
			}
			definition.Transitions = append(definition.Transitions, dto.ProjectManagementWorkflowTransitionDto{Id: item.Id, Name: item.Name, FromStageId: definition.Stages[i].Id, ToStageId: item.ToStageId, AllowedActorTypes: item.AllowedActorTypes})
		}
	}
	return definition, nil
}

func validateProjectManagementWorkflowDefinition(workflow dto.ProjectManagementWorkflowDto) error {
	stageIDs := make(map[int64]struct{}, len(workflow.Stages))
	pending, finished := 0, 0
	minSortOrder := -1
	for _, stage := range workflow.Stages {
		if stage.Id == 0 || strings.TrimSpace(stage.Name) == "" {
			return errorx.New("workflow stage id and name are required")
		}
		if !isProjectManagementStagePhase(stage.Phase) {
			return errorx.New("invalid workflow stage phase")
		}
		if _, exists := stageIDs[stage.Id]; exists {
			return errorx.New("workflow stage id must be unique")
		}
		stageIDs[stage.Id] = struct{}{}
		if minSortOrder < 0 || stage.SortOrder < minSortOrder {
			minSortOrder = stage.SortOrder
		}
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			pending++
		}
		if stage.Phase == constants.ProjectManagementStagePhaseFinished {
			finished++
		}
	}
	if pending != 1 {
		return errorx.New("workflow must have exactly one pending stage")
	}
	for _, stage := range workflow.Stages {
		if stage.Phase == constants.ProjectManagementStagePhasePending && stage.SortOrder != minSortOrder {
			return errorx.New("first workflow stage must be pending")
		}
	}
	if finished == 0 {
		return errorx.New("workflow must have at least one finished stage")
	}
	for _, stage := range workflow.Stages {
		if stage.SortOrder == maxWorkflowStageSortOrder(workflow.Stages) && stage.Phase != constants.ProjectManagementStagePhaseFinished {
			return errorx.New("last workflow stage must be finished")
		}
	}
	edges := make(map[string]struct{}, len(workflow.Transitions))
	for _, transition := range workflow.Transitions {
		if err := validateProjectManagementWorkflowTransition(transition); err != nil {
			return err
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
	}
	return nil
}

func maxWorkflowStageSortOrder(stages []dto.ProjectManagementWorkflowStageDto) int {
	maxOrder := stages[0].SortOrder
	for _, stage := range stages[1:] {
		maxOrder = max(maxOrder, stage.SortOrder)
	}
	return maxOrder
}

func (l *ProjectManagementWorkflowLogic) Detail(ctx context.Context, organizationId, id int64) (*model.ProjectManagementWorkflow, error) {
	workflow, err := l.repos.ProjectManagementWorkflowRepo.FindWorkflowSourceByID(ctx, id, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "workflow not found")
	}
	return workflow, nil
}

func (l *ProjectManagementWorkflowLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementWorkflowDto) (workflow *model.ProjectManagementWorkflow, err error) {
	if err = ValidateProjectManagementWorkflow(input); err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		template := &model.ProjectManagementTemplate{TemplateType: constants.ProjectManagementTemplateTypeOrganization, OrganizationId: claim.OrganizationId, Name: input.Name, Description: input.Description, Status: "Active"}
		if err := tx.WithContext(ctx).Create(template).Error; err != nil {
			return errorx.WithStack(err)
		}
		workflow = &model.ProjectManagementWorkflow{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: template.Id, Name: input.Name, Description: input.Description}
		if err = tx.WithContext(ctx).Omit("Kind", "OrganizationId", "ProjectId", "SourceTemplateId", "SourceTemplateName").Create(workflow).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = l.replaceWorkflowSourceDefinition(ctx, tx, workflow, input); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflow.Id, action: constants.OperationActionProjectManagementWorkflowCreate, after: workflowOperationValues(workflow)})
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
	if err = ValidateProjectManagementWorkflow(definition); err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflow, err = l.repos.ProjectManagementWorkflowRepo.LockWorkflowSourceWithTx(ctx, tx, id, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow not found")
		}
		before := workflowOperationValues(workflow)
		workflow.Name, workflow.Description = input.Name, input.Description
		workflow.RemoveRelationships()
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("id = ?", workflow.Id).Updates(workflow).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementTemplate{}).Where("id = ?", workflow.EntityId).Updates(map[string]any{"name": input.Name, "description": input.Description}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = l.replaceWorkflowSourceDefinition(ctx, tx, workflow, definition); err != nil {
			return err
		}
		old, next := operationValuesDiff(before, workflowOperationValues(workflow))
		if len(old) == 0 {
			return nil
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: workflow.Id, action: constants.OperationActionProjectManagementWorkflowUpdate, before: old, after: next})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim.OrganizationId, id)
}

func (l *ProjectManagementWorkflowLogic) replaceWorkflowSourceDefinition(ctx context.Context, tx *gorm.DB, workflow *model.ProjectManagementWorkflow, definition dto.ProjectManagementWorkflowDto) error {
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflow.Id).Delete(&model.ProjectManagementWorkflowStageTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflow.Id).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, workflow.EntityId).Delete(&model.ProjectManagementStage{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	stageIDs := make(map[int64]int64, len(definition.Stages))
	for _, input := range definition.Stages {
		stage := &model.ProjectManagementStage{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: workflow.EntityId, Name: input.Name, Description: input.Description, Phase: input.Phase, Color: input.Color, SortOrder: input.SortOrder}
		if err := tx.WithContext(ctx).Create(stage).Error; err != nil {
			return errorx.WithStack(err)
		}
		link := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, StageId: stage.Id}
		if err := tx.WithContext(ctx).Create(link).Error; err != nil {
			return errorx.WithStack(err)
		}
		stageIDs[input.Id] = stage.Id
	}
	for _, input := range definition.Transitions {
		transition := &model.ProjectManagementWorkflowStageTransition{WorkflowId: workflow.Id, Name: input.Name, FromStageId: stageIDs[input.FromStageId], ToStageId: stageIDs[input.ToStageId], AllowedActorTypes: input.AllowedActorTypes}
		if transition.FromStageId == 0 || transition.ToStageId == 0 {
			return errorx.New("workflow transition endpoint is invalid")
		}
		if err := tx.WithContext(ctx).Create(transition).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

func (l *ProjectManagementWorkflowLogic) SetDefault(ctx context.Context, claim utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflows, err := l.repos.ProjectManagementWorkflowRepo.LockWorkflowSourcesWithTx(ctx, tx, claim.OrganizationId)
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
			if err = tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("id IN ?", workflowIDs(workflows)).Update("is_default", false).Error; err != nil {
				return errorx.WithStack(err)
			}
			if err = tx.WithContext(ctx).Model(&model.ProjectManagementTemplate{}).Where("organization_id = ? AND template_type = ? AND (project_type IS NULL OR project_type = '')", claim.OrganizationId, constants.ProjectManagementTemplateTypeOrganization).Update("is_default", false).Error; err != nil {
				return errorx.WithStack(err)
			}
			if err = tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("id = ?", id).Update("is_default", true).Error; err != nil {
				return errorx.WithStack(err)
			}
			return errorx.WithStack(tx.WithContext(ctx).Model(&model.ProjectManagementTemplate{}).Where("id = ?", workflow.EntityId).Update("is_default", true).Error)
		}
		return errorx.New("workflow not found")
	})
}

func workflowIDs(workflows []*model.ProjectManagementWorkflow) []int64 {
	ids := make([]int64, 0, len(workflows))
	for _, workflow := range workflows {
		ids = append(ids, workflow.Id)
	}
	return ids
}

func (l *ProjectManagementWorkflowLogic) Delete(ctx context.Context, claim utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflow, err := l.repos.ProjectManagementWorkflowRepo.LockWorkflowSourceWithTx(ctx, tx, id, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "workflow not found")
		}
		if workflow.IsDefault {
			return errorx.New("default workflow cannot be deleted")
		}
		if err = l.repos.ProjectManagementWorkflowRepo.DeleteWorkflowSourceWithTx(ctx, tx, workflow); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementWorkflowOperation{claim: claim, workflowId: id, action: constants.OperationActionProjectManagementWorkflowDelete, before: workflowOperationValues(workflow)})
	})
}

func (l *ProjectManagementWorkflowLogic) recordOperation(ctx context.Context, tx *gorm.DB, input projectManagementWorkflowOperation) error {
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{WorkflowID: input.workflowId}, input.before, input.after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: input.claim.OrganizationId, EntityType: constants.OperationEntityTypeProjectManagementOrganization, OperatorId: input.claim.OrganizationUserId, OperatorName: input.claim.Username, OperatorType: constants.UserTypeBorrower, Operation: input.action, Metadata: metadata, OrganizationId: input.claim.OrganizationId})
}

func workflowOperationValues(workflow *model.ProjectManagementWorkflow) map[string]any {
	return map[string]any{"name": workflow.Name, "description": workflow.Description, "is_default": workflow.IsDefault}
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

func isProjectManagementStagePhase(phase constants.ProjectManagementStagePhase) bool {
	return phase == constants.ProjectManagementStagePhasePending || phase == constants.ProjectManagementStagePhaseInProgress || phase == constants.ProjectManagementStagePhaseFinished
}

func isProjectManagementMemberRole(role constants.ProjectManagementMemberRole) bool {
	return role == constants.ProjectManagementMemberRoleProjectOwner || role == constants.ProjectManagementMemberRoleProjectAdmin || role == constants.ProjectManagementMemberRoleTaskCreator || role == constants.ProjectManagementMemberRoleTaskExecutor
}
