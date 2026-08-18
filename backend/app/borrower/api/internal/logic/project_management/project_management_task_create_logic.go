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

type ProjectManagementTaskCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create project task
func NewProjectManagementTaskCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskCreateLogic {
	return &ProjectManagementTaskCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskCreateLogic) ProjectManagementTaskCreate(req *types.ProjectManagementTaskCreateRequest) (resp *types.ProjectManagementTaskEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskCreateRequest, borrower.ProjectManagementTaskCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskCreate(l.ctx, &borrower.ProjectManagementTaskCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
}
