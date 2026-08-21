package project_management_task

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

type AssigneeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// replace project task assignees
func NewAssigneeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssigneeUpdateLogic {
	return &AssigneeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssigneeUpdateLogic) AssigneeUpdate(req *types.ProjectManagementTaskAssigneeUpdateRequest) (resp *types.ProjectManagementTaskDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskAssigneeUpdateRequest, borrower.ProjectManagementTaskAssigneeUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.UpdateAssignees(l.ctx, &borrower.ProjectManagementTaskAssigneeUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskDetailResponse{ProjectManagementTaskEntity: *entity}, nil
}
