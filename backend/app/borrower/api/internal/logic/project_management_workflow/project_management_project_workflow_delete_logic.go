package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementProjectWorkflowDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete project workflow
func NewProjectManagementProjectWorkflowDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowDeleteLogic {
	return &ProjectManagementProjectWorkflowDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowDeleteLogic) ProjectManagementProjectWorkflowDelete(req *types.ProjectManagementProjectWorkflowDeleteRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowDeleteRequest, borrower.ProjectManagementProjectWorkflowDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.ProjectManagementWorkflowRpcClient.ProjectWorkflowDelete(l.ctx, &borrower.ProjectManagementProjectWorkflowDeleteRequest{Data: data}); err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
