package project_management_project

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

type WorkflowUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// workflow update
func NewWorkflowUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowUpdateLogic {
	return &WorkflowUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowUpdateLogic) WorkflowUpdate(req *types.ProjectManagementProjectWorkflowUpdateRequest) (resp *types.ProjectManagementWorkflowDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowUpdateRequest, borrower.ProjectManagementWorkflowUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Update(l.ctx, &borrower.ProjectManagementWorkflowUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowDetailResponse{ProjectManagementWorkflowEntity: *entity}, nil
}
