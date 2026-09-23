package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementProjectWorkflowStageUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update project workflow stages
func NewProjectManagementProjectWorkflowStageUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowStageUpdateLogic {
	return &ProjectManagementProjectWorkflowStageUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowStageUpdateLogic) ProjectManagementProjectWorkflowStageUpdate(req *types.ProjectManagementProjectWorkflowStageUpdateRequest) (resp *types.ProjectManagementWorkflowDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowStageUpdateRequest, borrower.ProjectManagementProjectWorkflowStageUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.ProjectWorkflowStageUpdate(l.ctx, &borrower.ProjectManagementProjectWorkflowStageUpdateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowDetailResponse{ProjectManagementWorkflowEntity: *entity}, nil
}
