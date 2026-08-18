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

type ProjectManagementWorkflowPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project workflows
func NewProjectManagementWorkflowPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowPageLogic {
	return &ProjectManagementWorkflowPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowPageLogic) ProjectManagementWorkflowPage(req *types.ProjectManagementWorkflowPageRequest) (resp *types.ProjectManagementWorkflowPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowPageRequest, borrower.ProjectManagementWorkflowPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.WorkflowPage(l.ctx, &borrower.ProjectManagementWorkflowPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementWorkflowPageResponseData, types.ProjectManagementWorkflowPageResponse](rpcResp.Data)
}
