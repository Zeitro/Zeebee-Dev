package projectmanagementtasktypelogic

import (
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func taskTypeEntity(taskType *model.ProjectManagementTaskType) (*proto.ProjectManagementTaskTypeEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementTaskType, proto.ProjectManagementTaskTypeEntity](taskType)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(taskType.CreatedAt), formatEntityTime(taskType.UpdatedAt)
	return entity, nil
}
