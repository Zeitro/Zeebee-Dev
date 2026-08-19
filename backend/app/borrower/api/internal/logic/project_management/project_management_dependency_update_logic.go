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

type ProjectManagementDependencyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update task dependency
func NewProjectManagementDependencyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementDependencyUpdateLogic {
	return &ProjectManagementDependencyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementDependencyUpdateLogic) ProjectManagementDependencyUpdate(req *types.ProjectManagementDependencyUpdateRequest) (resp *types.ProjectManagementTaskDependencyEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementDependencyUpdateRequest, borrower.ProjectManagementDependencyUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.DependencyUpdate(l.ctx, &borrower.ProjectManagementDependencyUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskDependencyEntity, types.ProjectManagementTaskDependencyEntity](rpcResp.Data)
}
