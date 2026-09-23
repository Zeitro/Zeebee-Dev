package project_management_task

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ViewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// read-only task view
func NewViewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ViewLogic {
	return &ViewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ViewLogic) View(req *types.ProjectManagementTaskViewRequest) (resp *types.ProjectManagementTaskViewEntity, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.View(l.ctx, &borrower.ProjectManagementTaskViewRequest{Data: &borrower.ProjectManagementTaskViewRequestData{Id: req.Id}})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementTaskViewEntity, types.ProjectManagementTaskViewEntity](rpcResp)
}
