package projectmanagementprojectlogic

import (
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func projectEntity(project *model.ProjectManagementProject) (*proto.ProjectManagementProjectEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementProject, proto.ProjectManagementProjectEntity](project)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(project.CreatedAt), formatEntityTime(project.UpdatedAt)
	return entity, nil
}
