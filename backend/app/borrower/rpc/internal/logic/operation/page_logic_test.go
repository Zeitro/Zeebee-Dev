package operationlogic

import (
	"testing"
	"time"

	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func TestOperationEntityReturnsMetadataWithoutProjectFields(t *testing.T) {
	createdAt := time.Date(2026, 8, 25, 7, 0, 0, 0, time.UTC)
	operation := &model.Operation{
		EntityId:   42,
		EntityType: constants.OperationEntityTypeProjectManagementProject,
		Operation:  constants.OperationActionProjectManagementTaskCreate,
		Metadata:   datatypes.JSON(`{"task_id":7,"title":"task"}`),
	}
	operation.CreatedAt = createdAt

	entity, err := operationEntity(operation)
	if err != nil {
		t.Fatal(err)
	}
	if entity.EntityId != operation.EntityId || entity.EntityType != string(operation.EntityType) || entity.Metadata != string(operation.Metadata) || entity.CreatedAt != createdAt.Format(time.RFC3339) {
		t.Fatalf("unexpected operation entity: %#v", entity)
	}
}
