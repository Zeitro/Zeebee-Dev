package project_management_task

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

type TransitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// transition task
func NewTransitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransitionLogic {
	return &TransitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransitionLogic) Transition(req *types.ProjectManagementTaskTransitionRequest) (resp *types.ProjectManagementTaskEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTransitionRequest, borrower.ProjectManagementTaskTransitionRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.Transition(l.ctx, &borrower.ProjectManagementTaskTransitionRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[proto.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
}
