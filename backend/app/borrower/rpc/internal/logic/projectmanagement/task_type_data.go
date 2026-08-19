package projectmanagementlogic

import (
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/model"
)

func taskTypeEntity(taskType *model.ProjectManagementTaskType) (*borrower.ProjectManagementTaskTypeEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementTaskType, borrower.ProjectManagementTaskTypeEntity](taskType)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(taskType.CreatedAt), formatEntityTime(taskType.UpdatedAt)
	return entity, nil
}
