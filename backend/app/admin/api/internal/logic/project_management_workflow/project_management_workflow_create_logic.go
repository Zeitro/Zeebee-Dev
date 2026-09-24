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

type ProjectManagementWorkflowCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewProjectManagementWorkflowCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowCreateLogic {
	return &ProjectManagementWorkflowCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowCreateLogic) ProjectManagementWorkflowCreate(req *types.ProjectManagementWorkflowCreateRequest) (resp *types.ProjectManagementWorkflowEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowCreateRequest, admin.ProjectManagementWorkflowCreateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Create(l.ctx, &admin.ProjectManagementWorkflowCreateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
}
