package logic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestValidateProjectManagementWorkflow(t *testing.T) {
	workflow := dto.ProjectManagementWorkflowDto{Name: "Source", Stages: []dto.ProjectManagementWorkflowStageDto{
		{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
		{Id: 2, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
	}, Transitions: []dto.ProjectManagementWorkflowTransitionDto{{Name: "Finish", FromStageId: 1, ToStageId: 2}}}
	require.NoError(t, ValidateProjectManagementWorkflow(workflow))
	workflow.Transitions[0].ToStageId = 3
	require.Error(t, ValidateProjectManagementWorkflow(workflow))
	workflow.Transitions[0].ToStageId = 2
	workflow.Transitions = append(workflow.Transitions, dto.ProjectManagementWorkflowTransitionDto{Name: "Finish again", FromStageId: 1, ToStageId: 2})
	require.ErrorContains(t, ValidateProjectManagementWorkflow(workflow), "workflow transition must be unique")
	workflow.Transitions = workflow.Transitions[:1]
	workflow.Stages[1].Phase = constants.ProjectManagementStagePhaseInProgress
	workflow.Stages = append(workflow.Stages, dto.ProjectManagementWorkflowStageDto{Id: 3, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 1})
	require.Error(t, ValidateProjectManagementWorkflow(workflow))
}

func TestProjectManagementWorkflowUpdateDefinition(t *testing.T) {
	definition, err := projectManagementWorkflowUpdateDefinition(dto.ProjectManagementWorkflowUpdateDto{Stages: []dto.ProjectManagementWorkflowUpdateStageDto{
		{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1, Transitions: []dto.ProjectManagementWorkflowUpdateTransitionDto{{Name: "Finish", FromStageId: 1, ToStageId: 2}}},
		{Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
	}})
	require.NoError(t, err)
	require.Len(t, definition.Transitions, 1)
	require.NotZero(t, definition.Stages[1].Id)
}

func TestProjectWorkflowTransitionsBindEveryDistinctStage(t *testing.T) {
	stages := []*model.ProjectManagementStage{
		{Base: model.Base{Id: 1}, Name: "Pending"},
		{Base: model.Base{Id: 2}, Name: "In progress"},
		{Base: model.Base{Id: 3}, Name: "Finished"},
	}
	transitions := projectWorkflowTransitions(stages)
	require.Len(t, transitions, 6)
	for _, transition := range transitions {
		require.NotEqual(t, transition.FromStageId, transition.ToStageId)
		require.NotEmpty(t, transition.Name)
		require.Equal(t, defaultWorkflowActorTypes(), transition.AllowedActorTypes)
	}
	_, err := validateProjectWorkflowTransitions(append(transitions, transitions[0]))
	require.ErrorContains(t, err, "workflow transition must be unique")
}

func TestDefaultWorkflowDefinitionIsValid(t *testing.T) {
	workflow := defaultWorkflowDefinition()
	workflow.Name = "Default"
	require.NoError(t, ValidateProjectManagementWorkflow(workflow))
	require.Len(t, workflow.Stages, 3)
	require.Len(t, workflow.Transitions, 2)
}
