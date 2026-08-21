package project_management_task

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register project task attachment
func NewAttachmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentCreateLogic {
	return &AttachmentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentCreateLogic) AttachmentCreate(req *types.ProjectManagementAttachmentCreateRequest) (resp *types.ProjectManagementAttachmentDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementAttachmentCreateRequest, borrower.ProjectManagementAttachmentCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.AttachmentCreate(l.ctx, &borrower.ProjectManagementAttachmentCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := attachmentEntity(rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementAttachmentDetailResponse{ProjectManagementAttachmentEntity: *entity}, nil
}
