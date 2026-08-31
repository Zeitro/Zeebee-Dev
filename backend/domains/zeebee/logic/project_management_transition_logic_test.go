package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func setupProjectManagementTransitionTest(t *testing.T, dependencyRule constants.ProjectManagementDependencyRule, predecessorPhase constants.ProjectManagementStagePhase) (*ProjectManagementTaskLogic, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_tasks (id integer primary key, is_del integer not null default 0, project_id integer not null default 0, organization_id integer not null default 0, parent_task_id integer not null default 0, stage_id integer not null default 0, creator_id integer not null default 0, creator_name text not null default '', title text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, is_del integer not null default 0, workflow_id integer not null default 0, name text not null default '', phase text not null default '', color text not null default '', sort_order integer not null default 0)",
		"CREATE TABLE project_management_task_dependencies (id integer primary key, is_del integer not null default 0, task_id integer not null, dependency_task_id integer not null, dependency_type text not null default 'PreTask', dependency_rule text not null default 'FS')",
		"CREATE TABLE project_management_task_assignees (id integer primary key, is_del integer not null default 0, task_id integer not null, stage_id integer not null default 0, organization_user_id integer not null, organization_user_full_name text not null default '')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, stage := range []struct {
		id    int64
		phase constants.ProjectManagementStagePhase
	}{{100, constants.ProjectManagementStagePhasePending}, {101, constants.ProjectManagementStagePhaseInProgress}, {102, constants.ProjectManagementStagePhaseFinished}} {
		if err := db.Exec("INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (?, 1, ?, ?, ?)", stage.id, stage.phase, stage.phase, stage.id).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Exec("INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id) VALUES (10, 1, 1, 100), (20, 1, 1, 100)").Error; err != nil {
		t.Fatal(err)
	}
	if dependencyRule != "" {
		if err := db.Exec("INSERT INTO project_management_task_dependencies (id, task_id, dependency_task_id, dependency_type, dependency_rule) VALUES (1, 10, 20, 'PreTask', ?)", dependencyRule).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Exec("UPDATE project_management_workflow_stages SET phase = ? WHERE id = 100", predecessorPhase).Error; err != nil {
		t.Fatal(err)
	}
	// Restore task 10's current stage independently; predecessor phase is represented by stage 200.
	if err := db.Exec("INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (200, 1, 'predecessor', ?, 200)", predecessorPhase).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("UPDATE project_management_tasks SET stage_id = 200 WHERE id = 20").Error; err != nil {
		t.Fatal(err)
	}
	base := NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{})
	return NewProjectManagementTaskLogic(base), db
}

func TestTransitionDependencyRules(t *testing.T) {
	cases := []struct {
		name       string
		dependency constants.ProjectManagementDependencyRule
		target     constants.ProjectManagementStagePhase
		blocked    bool
	}{
		{"FS blocks in progress", constants.ProjectManagementDependencyRuleFS, constants.ProjectManagementStagePhaseInProgress, true},
		{"SS blocks in progress", constants.ProjectManagementDependencyRuleSS, constants.ProjectManagementStagePhaseInProgress, true},
		{"FF blocks finished", constants.ProjectManagementDependencyRuleFF, constants.ProjectManagementStagePhaseFinished, true},
		{"SF blocks finished", constants.ProjectManagementDependencyRuleSF, constants.ProjectManagementStagePhaseFinished, true},
		{"direct finish ignores FS", constants.ProjectManagementDependencyRuleFS, constants.ProjectManagementStagePhaseFinished, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logic, db := setupProjectManagementTransitionTest(t, tc.dependency, constants.ProjectManagementStagePhasePending)
			current := &model.ProjectManagementWorkflowStage{Base: model.Base{Id: 100}, Phase: constants.ProjectManagementStagePhasePending}
			target := &model.ProjectManagementWorkflowStage{Base: model.Base{Id: 101}, Phase: tc.target}
			if tc.target == constants.ProjectManagementStagePhaseFinished {
				target.Id = 102
			}
			err := db.Transaction(func(tx *gorm.DB) error {
				return logic.validateTransitionDependencies(&projectManagementTransitionContext{ctx: context.Background(), tx: tx, task: &model.ProjectManagementTask{Base: model.Base{Id: 10}}, currentStage: current, targetStage: target})
			})
			if tc.blocked && err == nil {
				t.Fatal("blocked transition was accepted")
			}
			if !tc.blocked && err != nil {
				t.Fatalf("allowed direct finish rejected: %v", err)
			}
		})
	}
}

func TestTransitionActorAndChildRules(t *testing.T) {
	logic, db := setupProjectManagementTransitionTest(t, "", constants.ProjectManagementStagePhasePending)
	actor := utils.JwtUser{OrganizationUserId: 7}
	task := &model.ProjectManagementTask{Base: model.Base{Id: 10}}
	project := &model.ProjectManagementProject{OwnerId: 1, OwnerName: "Owner"}
	member := &model.ProjectManagementProjectMember{Role: constants.ProjectManagementMemberRoleProjectAdmin}
	transition := &model.ProjectManagementWorkflowTransition{AllowedActorTypes: []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleTaskExecutor}}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return logic.requireTransitionActor(&projectManagementTransitionContext{ctx: context.Background(), tx: tx, task: task, project: project, member: member, claim: actor, transition: transition})
	}); err == nil {
		t.Fatal("non-executor passed executor-only transition")
	}
	member.Role = constants.ProjectManagementMemberRoleTaskExecutor
	if err := db.Transaction(func(tx *gorm.DB) error {
		return logic.requireTransitionActor(&projectManagementTransitionContext{ctx: context.Background(), tx: tx, task: task, project: project, member: member, claim: actor, transition: transition})
	}); err != nil {
		t.Fatalf("executor member rejected: %v", err)
	}
	if err := db.Exec("INSERT INTO project_management_tasks (id, project_id, organization_id, parent_task_id, stage_id) VALUES (11, 1, 1, 10, 100)").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return logic.validateFinishedChildren(&projectManagementTransitionContext{ctx: context.Background(), tx: tx, task: task})
	}); err == nil {
		t.Fatal("unfinished child did not block finish")
	}
	if err := db.Exec("UPDATE project_management_tasks SET stage_id = 102 WHERE id = 11").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return logic.validateFinishedChildren(&projectManagementTransitionContext{ctx: context.Background(), tx: tx, task: task})
	}); err != nil {
		t.Fatalf("finished child blocked completion: %v", err)
	}
}
