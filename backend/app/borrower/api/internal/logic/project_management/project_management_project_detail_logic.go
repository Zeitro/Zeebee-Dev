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

type ProjectManagementProjectDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get project
func NewProjectManagementProjectDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectDetailLogic {
	return &ProjectManagementProjectDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectDetailLogic) ProjectManagementProjectDetail(req *types.ProjectManagementProjectDetailRequest) (resp *types.ProjectManagementProjectEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectDetailRequest, borrower.ProjectManagementProjectDetailRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.ProjectDetail(l.ctx, &borrower.ProjectManagementProjectDetailRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementProjectEntity, types.ProjectManagementProjectEntity](rpcResp.Data)
}
