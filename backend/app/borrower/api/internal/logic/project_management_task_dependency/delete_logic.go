package project_management_task_dependency

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete task dependency
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.ProjectManagementDependencyDeleteRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementDependencyDeleteRequest, borrower.ProjectManagementDependencyDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProjectManagementTaskDependencyRpcClient.Delete(l.ctx, &borrower.ProjectManagementDependencyDeleteRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
