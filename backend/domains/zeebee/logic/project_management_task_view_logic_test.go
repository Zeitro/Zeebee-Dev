package logic

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestTaskViewTimelineUsesTaskAndProjectDateUnion(t *testing.T) {
	view := &dto.ProjectManagementTaskViewDto{}
	updateTaskTimeline(view, &sql.NullTime{Time: time.Date(2026, time.March, 10, 0, 0, 0, 0, time.UTC), Valid: true})
	updateTaskTimeline(view, &sql.NullTime{Time: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC), Valid: true})
	updateTaskTimeline(view, &sql.NullTime{Time: time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC), Valid: true})
	require.Equal(t, "2026-03-01", view.TimelineStartAt.Time.Format(time.DateOnly))
	require.Equal(t, "2026-03-20", view.TimelineEndAt.Time.Format(time.DateOnly))
}

func TestTaskOutOfPlan(t *testing.T) {
	project := &model.ProjectManagementProject{
		ProjectType: constants.ProjectManagementProjectTypePlanned,
		StartAt:     &sql.NullTime{Time: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndAt:       &sql.NullTime{Time: time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC), Valid: true},
	}
	require.False(t, taskOutOfPlan(project, &sql.NullTime{Time: time.Date(2026, time.March, 31, 23, 0, 0, 0, time.UTC), Valid: true}, nil))
	require.True(t, taskOutOfPlan(project, nil, &sql.NullTime{Time: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), Valid: true}))
}
