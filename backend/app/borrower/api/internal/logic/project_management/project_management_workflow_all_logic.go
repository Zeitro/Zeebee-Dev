package project_management

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementWorkflowAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// list enabled project workflows
func NewProjectManagementWorkflowAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowAllLogic {
	return &ProjectManagementWorkflowAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowAllLogic) ProjectManagementWorkflowAll(req *types.EmptyRequest) (resp *types.ProjectManagementWorkflowAllResponse, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.WorkflowAll(l.ctx, &borrower.ProjectManagementWorkflowAllRequest{
		Data: &borrower.ProjectManagementWorkflowAllRequestData{},
	})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*borrower.ProjectManagementWorkflowEntity, []*types.ProjectManagementWorkflowEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowAllResponse{Data: *data}, nil

}
