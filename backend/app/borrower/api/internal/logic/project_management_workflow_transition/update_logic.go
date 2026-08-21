package project_management_workflow_transition

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

// update
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProjectManagementWorkflowTransitionUpdateRequest) (resp *types.ProjectManagementWorkflowTransitionDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowTransitionUpdateRequest, borrower.ProjectManagementWorkflowTransitionUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.WorkflowTransitionRpcClient.Update(l.ctx, &borrower.ProjectManagementWorkflowTransitionUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowTransitionEntity, types.ProjectManagementWorkflowTransitionEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowTransitionDetailResponse{ProjectManagementWorkflowTransitionEntity: *entity}, nil
}
