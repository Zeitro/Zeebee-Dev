package projectmanagementlogic

import (
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/model"
)

func projectEntity(project *model.ProjectManagementProject) (*borrower.ProjectManagementProjectEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementProject, borrower.ProjectManagementProjectEntity](project)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(project.CreatedAt), formatEntityTime(project.UpdatedAt)
	return entity, nil
}
