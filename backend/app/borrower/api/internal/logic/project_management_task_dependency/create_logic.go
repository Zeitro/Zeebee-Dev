package project_management_task_dependency

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

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create task dependency
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProjectManagementDependencyCreateRequest) (resp *types.ProjectManagementTaskDependencyEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementDependencyCreateRequest, borrower.ProjectManagementDependencyCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.TaskId
	rpcResp, err := l.svcCtx.ProjectManagementTaskDependencyRpcClient.Create(l.ctx, &borrower.ProjectManagementDependencyCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[proto.ProjectManagementTaskDependencyEntity, types.ProjectManagementTaskDependencyEntity](rpcResp.Data)
}
