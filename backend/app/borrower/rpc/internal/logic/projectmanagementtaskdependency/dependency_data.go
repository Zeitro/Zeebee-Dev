package projectmanagementtaskdependencylogic

import (
	"time"

	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func dependencyEntity(dependency *model.ProjectManagementTaskDependency) (*proto.ProjectManagementTaskDependencyEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementTaskDependency, proto.ProjectManagementTaskDependencyEntity](dependency)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(dependency.CreatedAt), formatEntityTime(dependency.UpdatedAt)
	return entity, nil
}

func formatEntityTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}
