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

type ProjectManagementTaskTransitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// transition task
func NewProjectManagementTaskTransitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskTransitionLogic {
	return &ProjectManagementTaskTransitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskTransitionLogic) ProjectManagementTaskTransition(req *types.ProjectManagementTaskTransitionRequest) (resp *types.ProjectManagementTaskEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTransitionRequest, borrower.ProjectManagementTaskTransitionRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskTransition(l.ctx, &borrower.ProjectManagementTaskTransitionRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
}
