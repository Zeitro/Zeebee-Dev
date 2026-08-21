package project_management_workflow_stage

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.ProjectManagementWorkflowStageDetailRequest) (resp *types.EmptyResponse, err error) {
	if _, err = l.svcCtx.WorkflowStageRpcClient.Delete(l.ctx, &borrower.ProjectManagementWorkflowStageDetailRequest{Data: &borrower.ProjectManagementWorkflowStageDetailRequestData{WorkflowId: req.WorkflowId, Id: req.Id}}); err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
