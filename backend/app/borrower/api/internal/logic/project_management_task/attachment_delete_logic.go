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

type AttachmentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete project task attachment
func NewAttachmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDeleteLogic {
	return &AttachmentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentDeleteLogic) AttachmentDelete(req *types.ProjectManagementAttachmentDeleteRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementAttachmentDeleteRequest, borrower.ProjectManagementAttachmentDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	data.TaskId = req.Id
	_, err = l.svcCtx.ProjectManagementTaskRpcClient.AttachmentDelete(l.ctx, &borrower.ProjectManagementAttachmentDeleteRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
