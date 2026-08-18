package project_management

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementAttachmentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register project task attachment
func NewProjectManagementAttachmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementAttachmentCreateLogic {
	return &ProjectManagementAttachmentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementAttachmentCreateLogic) ProjectManagementAttachmentCreate(req *types.ProjectManagementAttachmentCreateRequest) (resp *types.ProjectManagementAttachmentEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementAttachmentCreateRequest, borrower.ProjectManagementAttachmentCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.AttachmentCreate(l.ctx, &borrower.ProjectManagementAttachmentCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementAttachmentEntity, types.ProjectManagementAttachmentEntity](rpcResp.Data)
}
