package project_management_project

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkflowDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// workflow detail
func NewWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowDetailLogic {
	return &WorkflowDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowDetailLogic) WorkflowDetail(req *types.ProjectManagementProjectWorkflowDetailRequest) (resp *types.ProjectManagementWorkflowEntity, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Detail(l.ctx, &borrower.ProjectManagementWorkflowDetailRequest{Data: &borrower.ProjectManagementWorkflowDetailRequestData{ProjectId: req.Id}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
}
