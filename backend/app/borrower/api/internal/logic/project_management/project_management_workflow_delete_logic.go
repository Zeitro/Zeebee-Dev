package project_management

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementWorkflowDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete project workflow
func NewProjectManagementWorkflowDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowDeleteLogic {
	return &ProjectManagementWorkflowDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowDeleteLogic) ProjectManagementWorkflowDelete(req *types.ProjectManagementWorkflowDeleteRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowDeleteRequest, borrower.ProjectManagementWorkflowDeleteRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if _, err = l.svcCtx.ProjectManagementRpcClient.WorkflowDelete(l.ctx, &borrower.ProjectManagementWorkflowDeleteRequest{Data: data}); err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
