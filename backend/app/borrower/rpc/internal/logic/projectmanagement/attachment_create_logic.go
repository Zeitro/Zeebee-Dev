package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type AttachmentCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentCreateLogic {
	return &AttachmentCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentCreateLogic) AttachmentCreate(in *borrower.ProjectManagementAttachmentCreateRequest) (resp *borrower.ProjectManagementAttachmentDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	file, err := base.ProjectManagementLogic.RegisterTaskAttachment(l.ctx, actor, dto.ProjectManagementAttachmentCreateDto{TaskId: in.Data.TaskId, FileName: in.Data.FileName, FileUrl: in.Data.FileUrl, StorageType: constants.StorageType(in.Data.StorageType)})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := attachmentEntity(file)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementAttachmentDetailResponse{Data: data}, nil
}
