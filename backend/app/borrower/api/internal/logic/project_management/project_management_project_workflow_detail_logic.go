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

type ProjectManagementProjectWorkflowDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get project workflow
func NewProjectManagementProjectWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowDetailLogic {
	return &ProjectManagementProjectWorkflowDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowDetailLogic) ProjectManagementProjectWorkflowDetail(req *types.ProjectManagementProjectWorkflowDetailRequest) (resp *types.ProjectManagementWorkflowEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowDetailRequest, borrower.ProjectManagementProjectWorkflowDetailRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.ProjectWorkflowDetail(l.ctx, &borrower.ProjectManagementProjectWorkflowDetailRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
}
