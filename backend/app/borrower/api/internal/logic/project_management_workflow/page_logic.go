package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project workflows
func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(req *types.ProjectManagementWorkflowPageRequest) (resp *types.ProjectManagementWorkflowPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowPageRequest, borrower.ProjectManagementWorkflowPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Page(l.ctx, &borrower.ProjectManagementWorkflowPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementWorkflowPageResponseData, types.ProjectManagementWorkflowPageResponse](rpcResp.Data)
}
