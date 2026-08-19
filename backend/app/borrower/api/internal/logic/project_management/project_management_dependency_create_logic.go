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

type ProjectManagementDependencyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create task dependency
func NewProjectManagementDependencyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementDependencyCreateLogic {
	return &ProjectManagementDependencyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementDependencyCreateLogic) ProjectManagementDependencyCreate(req *types.ProjectManagementDependencyCreateRequest) (resp *types.ProjectManagementTaskDependencyEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementDependencyCreateRequest, borrower.ProjectManagementDependencyCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.DependencyCreate(l.ctx, &borrower.ProjectManagementDependencyCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskDependencyEntity, types.ProjectManagementTaskDependencyEntity](rpcResp.Data)
}
