package projectmanagementlogic

import (
	"testing"
	"time"

	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/model"
)

func TestOperationEntityKeepsProjectTaskContext(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	updatedAt := now.Add(time.Hour)
	operation := &model.Operation{
		Base:         model.Base{Id: 1, CreatedAt: now, UpdatedAt: updatedAt},
		EntityId:     2,
		OperatorId:   3,
		OperatorName: "Operator",
		Operation:    "task_update",
		Metadata:     datatypes.JSON([]byte(`{"task_id":4}`)),
	}
	entity, err := operationEntity(operation)
	if err != nil {
		t.Fatal(err)
	}
	if entity.ProjectId != 2 || entity.TaskId != 4 || entity.OperatorId != 3 || entity.CreatedAt != now.Format(time.RFC3339) || entity.UpdatedAt != updatedAt.Format(time.RFC3339) {
		t.Fatalf("operation context lost: %#v", entity)
	}
}
