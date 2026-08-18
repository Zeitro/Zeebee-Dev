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

type ProjectManagementTaskAssigneeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// replace project task assignees
func NewProjectManagementTaskAssigneeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskAssigneeUpdateLogic {
	return &ProjectManagementTaskAssigneeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskAssigneeUpdateLogic) ProjectManagementTaskAssigneeUpdate(req *types.ProjectManagementTaskAssigneeUpdateRequest) (resp *types.ProjectManagementTaskEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskAssigneeUpdateRequest, borrower.ProjectManagementTaskAssigneeUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskAssigneeUpdate(l.ctx, &borrower.ProjectManagementTaskAssigneeUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
}
