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

type ProjectManagementProjectCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create project
func NewProjectManagementProjectCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectCreateLogic {
	return &ProjectManagementProjectCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectCreateLogic) ProjectManagementProjectCreate(req *types.ProjectManagementProjectCreateRequest) (resp *types.ProjectManagementProjectEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectCreateRequest, borrower.ProjectManagementProjectCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.ProjectCreate(l.ctx, &borrower.ProjectManagementProjectCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementProjectEntity, types.ProjectManagementProjectEntity](rpcResp.Data)
}
