package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementProjectWorkflowDefaultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// set project workflow default
func NewProjectManagementProjectWorkflowDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowDefaultLogic {
	return &ProjectManagementProjectWorkflowDefaultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowDefaultLogic) ProjectManagementProjectWorkflowDefault(req *types.ProjectManagementProjectWorkflowDefaultRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowDefaultRequest, borrower.ProjectManagementProjectWorkflowDefaultRequestData](req)
	if err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.ProjectManagementWorkflowRpcClient.ProjectWorkflowDefault(l.ctx, &borrower.ProjectManagementProjectWorkflowDefaultRequest{Data: data}); err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
