package project_management_workflow

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementWorkflowUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewProjectManagementWorkflowUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowUpdateLogic {
	return &ProjectManagementWorkflowUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowUpdateLogic) ProjectManagementWorkflowUpdate(req *types.ProjectManagementWorkflowUpdateRequest) (resp *types.ProjectManagementWorkflowEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowUpdateRequest, admin.ProjectManagementWorkflowUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Update(l.ctx, &admin.ProjectManagementWorkflowUpdateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
}
