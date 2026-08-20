package project_management_task

import (
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

func attachmentEntity(file *proto.FileEntity) (*types.ProjectManagementAttachmentEntity, error) {
	entity, err := copierx.Trans[proto.FileEntity, types.ProjectManagementAttachmentEntity](file)
	if err != nil {
		return nil, err
	}
	entity.FileId = file.Id
	return entity, nil
}
