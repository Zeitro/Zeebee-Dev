package project_management_workflow_stage

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

type AllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLogic) All(req *types.ProjectManagementWorkflowStageAllRequest) (resp *types.ProjectManagementWorkflowStageAllResponse, err error) {
	rpcResp, err := l.svcCtx.WorkflowStageRpcClient.All(l.ctx, &borrower.ProjectManagementWorkflowStageAllRequest{Data: &borrower.ProjectManagementWorkflowStageAllRequestData{WorkflowId: req.WorkflowId}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := copierx.Trans[[]*proto.ProjectManagementWorkflowStageEntity, []*types.ProjectManagementWorkflowStageEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowStageAllResponse{Lists: *data}, nil
}
