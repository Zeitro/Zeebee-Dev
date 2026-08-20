package project_management_task_dependency

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update task dependency
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProjectManagementDependencyUpdateRequest) (resp *types.ProjectManagementTaskDependencyEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementDependencyUpdateRequest, borrower.ProjectManagementDependencyUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.TaskId
	rpcResp, err := l.svcCtx.ProjectManagementTaskDependencyRpcClient.Update(l.ctx, &borrower.ProjectManagementDependencyUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[proto.ProjectManagementTaskDependencyEntity, types.ProjectManagementTaskDependencyEntity](rpcResp.Data)
}
