package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if _, err = base.ProjectManagementTaskLogic.DetailForProject(l.ctx, claim, in.Data.ProjectId, in.Data.TaskId); err != nil {
		return nil, err
	}
	file, err := base.ProjectManagementAttachmentLogic.Create(l.ctx, claim, dto.ProjectManagementAttachmentCreateDto{TaskId: in.Data.TaskId, FileName: in.Data.FileName, FileUrl: in.Data.FileUrl, StorageType: constants.StorageType(in.Data.StorageType), UploadType: constants.UploadType(in.Data.UploadType)})
	if err != nil {
		return nil, err
	}
	data, err := attachmentEntity(file)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementAttachmentDetailResponse{Data: data}, nil
}

func attachmentEntity(file *model.File) (*proto.FileEntity, error) {
	entity, err := copierx.Trans[model.File, proto.FileEntity](file,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}
	entity.Id = file.Id
	return entity, nil
}
