package project_management_task

import (
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

func attachmentEntity(file *proto.FileEntity) (*types.FileEntity, error) {
	return copierx.Trans[proto.FileEntity, types.FileEntity](file)
}
