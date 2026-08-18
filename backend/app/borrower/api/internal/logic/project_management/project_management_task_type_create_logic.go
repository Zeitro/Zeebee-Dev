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

type ProjectManagementTaskTypeCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create project task type
func NewProjectManagementTaskTypeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskTypeCreateLogic {
	return &ProjectManagementTaskTypeCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskTypeCreateLogic) ProjectManagementTaskTypeCreate(req *types.ProjectManagementTaskTypeCreateRequest) (resp *types.ProjectManagementTaskTypeEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTypeCreateRequest, borrower.ProjectManagementTaskTypeCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskTypeCreate(l.ctx, &borrower.ProjectManagementTaskTypeCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskTypeEntity, types.ProjectManagementTaskTypeEntity](rpcResp.Data)
}
