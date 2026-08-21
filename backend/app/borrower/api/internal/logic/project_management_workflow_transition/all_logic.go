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

func (l *AllLogic) All(req *types.ProjectManagementWorkflowTransitionAllRequest) (resp *types.ProjectManagementWorkflowTransitionAllResponse, err error) {
	rpcResp, err := l.svcCtx.WorkflowTransitionRpcClient.All(l.ctx, &borrower.ProjectManagementWorkflowTransitionAllRequest{Data: &borrower.ProjectManagementWorkflowTransitionAllRequestData{WorkflowId: req.WorkflowId, FromStageId: req.FromStageId}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := copierx.Trans[[]*proto.ProjectManagementWorkflowTransitionEntity, []*types.ProjectManagementWorkflowTransitionEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowTransitionAllResponse{Lists: *data}, nil
}
