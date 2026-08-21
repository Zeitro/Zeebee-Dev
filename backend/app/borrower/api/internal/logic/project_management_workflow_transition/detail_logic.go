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

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.ProjectManagementWorkflowTransitionDetailRequest) (resp *types.ProjectManagementWorkflowTransitionDetailResponse, err error) {
	rpcResp, err := l.svcCtx.WorkflowTransitionRpcClient.Detail(l.ctx, &borrower.ProjectManagementWorkflowTransitionDetailRequest{Data: &borrower.ProjectManagementWorkflowTransitionDetailRequestData{WorkflowId: req.WorkflowId, Id: req.Id}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowTransitionEntity, types.ProjectManagementWorkflowTransitionEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowTransitionDetailResponse{ProjectManagementWorkflowTransitionEntity: *entity}, nil
}
