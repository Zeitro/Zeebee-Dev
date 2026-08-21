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

func (l *DetailLogic) Detail(req *types.ProjectManagementWorkflowStageDetailRequest) (resp *types.ProjectManagementWorkflowStageDetailResponse, err error) {
	rpcResp, err := l.svcCtx.WorkflowStageRpcClient.Detail(l.ctx, &borrower.ProjectManagementWorkflowStageDetailRequest{Data: &borrower.ProjectManagementWorkflowStageDetailRequestData{WorkflowId: req.WorkflowId, Id: req.Id}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowStageEntity, types.ProjectManagementWorkflowStageEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowStageDetailResponse{ProjectManagementWorkflowStageEntity: *entity}, nil
}
