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

type ProjectManagementTaskTypeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update project task type
func NewProjectManagementTaskTypeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskTypeUpdateLogic {
	return &ProjectManagementTaskTypeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskTypeUpdateLogic) ProjectManagementTaskTypeUpdate(req *types.ProjectManagementTaskTypeUpdateRequest) (resp *types.ProjectManagementTaskTypeEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTypeUpdateRequest, borrower.ProjectManagementTaskTypeUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskTypeUpdate(l.ctx, &borrower.ProjectManagementTaskTypeUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskTypeEntity, types.ProjectManagementTaskTypeEntity](rpcResp.Data)
}
