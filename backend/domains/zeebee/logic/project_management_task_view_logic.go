package logic

import (
	"context"
	"database/sql"
	"sort"
	"time"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func (l *ProjectManagementTaskLogic) View(ctx context.Context, claim utils.JwtUser, projectId int64) (*dto.ProjectManagementTaskViewDto, error) {
	project, err := l.projectManagementLogic().RequireProject(ctx, projectId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	if _, err = l.projectManagementLogic().RequireMember(ctx, projectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	tasks, err := l.repos.ProjectManagementTaskRepo.ListByProjectForView(ctx, projectId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	taskIds := make([]int64, 0, len(tasks))
	for _, task := range tasks {
		taskIds = append(taskIds, task.Id)
	}
	assignees, err := l.repos.ProjectManagementTaskRepo.ListCurrentAssigneesByTaskIds(ctx, taskIds)
	if err != nil {
		return nil, err
	}
	dependencies, err := l.repos.ProjectManagementDependencyRepo.ListByTaskIds(ctx, taskIds)
	if err != nil {
		return nil, err
	}
	result := &dto.ProjectManagementTaskViewDto{ProjectId: projectId, ProjectStartAt: project.StartAt, ProjectEndAt: project.EndAt, Tasks: make([]*dto.ProjectManagementTaskDto, 0, len(tasks))}
	updateTaskTimeline(result, project.StartAt)
	updateTaskTimeline(result, project.EndAt)
	byID := make(map[int64]*dto.ProjectManagementTaskDto, len(tasks))
	for _, task := range tasks {
		item, convertErr := copierx.Trans[model.ProjectManagementTask, dto.ProjectManagementTaskDto](task)
		if convertErr != nil {
			return nil, convertErr
		}
		item.Assignees = make([]*dto.ProjectManagementTaskAssigneeDto, 0)
		item.Children = make([]*dto.ProjectManagementTaskDto, 0)
		item.Dependencies = make([]*dto.ProjectManagementTaskDependencyDto, 0)
		item.OutOfPlan = taskOutOfPlan(project, item.StartAt, item.EndAt)
		byID[item.Id] = item
	}
	for _, assignee := range assignees {
		if item := byID[assignee.TaskId]; item != nil {
			converted, convertErr := copierx.Trans[model.ProjectManagementTaskAssignee, dto.ProjectManagementTaskAssigneeDto](assignee)
			if convertErr != nil {
				return nil, convertErr
			}
			item.Assignees = append(item.Assignees, converted)
		}
	}
	for _, dependency := range dependencies {
		item := byID[dependency.TaskId]
		if item == nil {
			continue
		}
		converted, convertErr := copierx.Trans[model.ProjectManagementTaskDependency, dto.ProjectManagementTaskDependencyDto](dependency)
		if convertErr != nil {
			return nil, convertErr
		}
		if dependency.DependencyTask != nil {
			converted.DependencyTaskTitle = dependency.DependencyTask.Title
		}
		item.Dependencies = append(item.Dependencies, converted)
	}
	for _, item := range byID {
		if item.ParentTaskId == 0 {
			result.Tasks = append(result.Tasks, item)
			continue
		}
		if parent := byID[item.ParentTaskId]; parent != nil {
			parent.Children = append(parent.Children, item)
		} else {
			result.Tasks = append(result.Tasks, item)
		}
	}
	sort.Slice(result.Tasks, func(i, j int) bool { return result.Tasks[i].Id < result.Tasks[j].Id })
	for _, item := range byID {
		sort.Slice(item.Children, func(i, j int) bool { return item.Children[i].Id < item.Children[j].Id })
	}
	for _, item := range byID {
		updateTaskTimeline(result, item.StartAt)
		updateTaskTimeline(result, item.EndAt)
	}
	if result.TimelineStartAt == nil {
		result.TimelineStartAt = result.ProjectStartAt
	}
	if result.TimelineEndAt == nil {
		result.TimelineEndAt = result.ProjectEndAt
	}
	return result, nil
}

func taskOutOfPlan(project *model.ProjectManagementProject, startAt, endAt *sql.NullTime) bool {
	if project.ProjectType != constants.ProjectManagementProjectTypePlanned || project.StartAt == nil || !project.StartAt.Valid || project.EndAt == nil || !project.EndAt.Valid {
		return false
	}
	for _, value := range []*sql.NullTime{startAt, endAt} {
		if value == nil || !value.Valid {
			continue
		}
		date := time.Date(value.Time.Year(), value.Time.Month(), value.Time.Day(), 0, 0, 0, 0, time.UTC)
		start := time.Date(project.StartAt.Time.Year(), project.StartAt.Time.Month(), project.StartAt.Time.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(project.EndAt.Time.Year(), project.EndAt.Time.Month(), project.EndAt.Time.Day(), 0, 0, 0, 0, time.UTC)
		if date.Before(start) || date.After(end) {
			return true
		}
	}
	return false
}

func updateTaskTimeline(view *dto.ProjectManagementTaskViewDto, value *sql.NullTime) {
	if value == nil || !value.Valid {
		return
	}
	if view.TimelineStartAt == nil || !view.TimelineStartAt.Valid || value.Time.Before(view.TimelineStartAt.Time) {
		view.TimelineStartAt = &sql.NullTime{Time: value.Time, Valid: true}
	}
	if view.TimelineEndAt == nil || !view.TimelineEndAt.Valid || value.Time.After(view.TimelineEndAt.Time) {
		view.TimelineEndAt = &sql.NullTime{Time: value.Time, Valid: true}
	}
}
