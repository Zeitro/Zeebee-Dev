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
		{"reopened", &sql.NullTime{Time: now.Add(-time.Hour), Valid: true}, constants.ProjectManagementStagePhasePending, true, int64(time.Hour.Seconds())},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			overdue, seconds := projectManagementTaskOverdue(&model.ProjectManagementTask{EndAt: tc.endAt}, tc.phase, now)
			if overdue != tc.overdue || seconds != tc.seconds {
				t.Fatalf("got overdue=%t seconds=%d", overdue, seconds)
			}
		})
	}
}
