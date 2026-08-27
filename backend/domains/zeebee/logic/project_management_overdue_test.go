package logic

import (
	"database/sql"
	"testing"
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectManagementTaskOverdue(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name    string
		endAt   *sql.NullTime
		phase   constants.ProjectManagementStagePhase
		overdue bool
		seconds int64
	}{
		{"no deadline", nil, constants.ProjectManagementStagePhaseInProgress, false, 0},
		{"not due", &sql.NullTime{Time: now.Add(time.Hour), Valid: true}, constants.ProjectManagementStagePhaseInProgress, false, 0},
		{"finished", &sql.NullTime{Time: now.Add(-time.Hour), Valid: true}, constants.ProjectManagementStagePhaseFinished, false, 0},
		{"reopened after deadline day", &sql.NullTime{Time: now.AddDate(0, 0, -1).Truncate(24 * time.Hour), Valid: true}, constants.ProjectManagementStagePhasePending, true, int64(12 * time.Hour.Seconds())},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			overdue, seconds := projectManagementTaskOverdue(&model.ProjectManagementTask{EndAt: tc.endAt}, tc.phase, now)
			if overdue != tc.overdue || seconds != tc.seconds {
				t.Fatalf("got overdue=%t seconds=%d", overdue, seconds)
			}
		})
	}
	endAt := &sql.NullTime{Time: time.Date(2026, time.May, 3, 0, 0, 0, 0, time.UTC), Valid: true}
	if overdue, _ := projectManagementTaskOverdue(&model.ProjectManagementTask{EndAt: endAt}, constants.ProjectManagementStagePhaseInProgress, time.Date(2026, time.May, 3, 23, 59, 59, 0, time.UTC)); overdue {
		t.Fatal("task became overdue before the day after its deadline")
	}
	if overdue, _ := projectManagementTaskOverdue(&model.ProjectManagementTask{EndAt: endAt}, constants.ProjectManagementStagePhaseInProgress, time.Date(2026, time.May, 4, 0, 0, 0, 0, time.UTC)); !overdue {
		t.Fatal("task was not overdue at the start of the day after its deadline")
	}
}
