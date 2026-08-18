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

type ProjectManagementTaskUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update project task
func NewProjectManagementTaskUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskUpdateLogic {
	return &ProjectManagementTaskUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskUpdateLogic) ProjectManagementTaskUpdate(req *types.ProjectManagementTaskUpdateRequest) (resp *types.ProjectManagementTaskEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskUpdateRequest, borrower.ProjectManagementTaskUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskUpdate(l.ctx, &borrower.ProjectManagementTaskUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
}
