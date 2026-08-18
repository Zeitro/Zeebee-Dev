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

type ProjectManagementProjectUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update project
func NewProjectManagementProjectUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectUpdateLogic {
	return &ProjectManagementProjectUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectUpdateLogic) ProjectManagementProjectUpdate(req *types.ProjectManagementProjectUpdateRequest) (resp *types.ProjectManagementProjectEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectUpdateRequest, borrower.ProjectManagementProjectUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.ProjectUpdate(l.ctx, &borrower.ProjectManagementProjectUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementProjectEntity, types.ProjectManagementProjectEntity](rpcResp.Data)
}
