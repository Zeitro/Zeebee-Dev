package repo_test

import (
	"context"
	"sync"
	"testing"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	domainlogic "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/stretchr/testify/require"
)

func TestProjectManagementDependencyPostgresOppositeEdges(t *testing.T) {
	db := projectManagementPostgresDB(t)
	project := &model.ProjectManagementProject{OrganizationId: 1, Name: "Renovation", OwnerId: 1, OwnerName: "Owner"}
	require.NoError(t, db.Create(project).Error)
	require.NoError(t, db.Create(&model.ProjectManagementProjectMember{
		ProjectId: project.Id, OrganizationUserId: 1, Role: constants.ProjectManagementMemberRoleTaskExecutor,
	}).Error)
	stage := &model.ProjectManagementWorkflowStage{WorkflowId: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10}
	require.NoError(t, db.Create(stage).Error)
	tasks := []*model.ProjectManagementTask{
		{ProjectId: project.Id, OrganizationId: 1, StageId: stage.Id, CreatorId: 1, Title: "First", TaskType: &dto.ProjectManagementTaskTypeSnapshotDto{}},
		{ProjectId: project.Id, OrganizationId: 1, StageId: stage.Id, CreatorId: 1, Title: "Second", TaskType: &dto.ProjectManagementTaskTypeSnapshotDto{}},
	}
	require.NoError(t, db.Create(&tasks).Error)

	projectLogic := domainlogic.NewProjectManagementDependencyLogic(
		domainlogic.NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}),
	)
	actor := &utils.JwtUser{OrganizationId: 1, OrganizationUserId: 1, Username: "Executor"}
	start := make(chan struct{})
	errs := make(chan error, 2)
	inputs := []dto.ProjectManagementTaskDependencyDto{
		{TaskId: tasks[0].Id, DependencyTaskId: tasks[1].Id, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS},
		{TaskId: tasks[1].Id, DependencyTaskId: tasks[0].Id, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS},
	}
	var wg sync.WaitGroup
	for _, input := range inputs {
		wg.Add(1)
		go func(input dto.ProjectManagementTaskDependencyDto) {
			defer wg.Done()
			<-start
			_, err := projectLogic.Create(context.Background(), actor, input)
			errs <- err
		}(input)
	}
	close(start)
	wg.Wait()
	close(errs)

	successes := 0
	for err := range errs {
		if err == nil {
			successes++
			continue
		}
		require.ErrorContains(t, err, "cycle")
	}
	require.Equal(t, 1, successes)
	var count int64
	require.NoError(t, db.Model(&model.ProjectManagementTaskDependency{}).Count(&count).Error)
	require.EqualValues(t, 1, count)
}
