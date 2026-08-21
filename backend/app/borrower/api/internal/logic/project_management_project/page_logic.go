package project_management_project

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page projects
func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(req *types.ProjectManagementProjectPageRequest) (resp *types.ProjectManagementProjectPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectPageRequest, borrower.ProjectManagementProjectPageRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementProjectRpcClient.Page(l.ctx, &borrower.ProjectManagementProjectPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementProjectPageResponseData, types.ProjectManagementProjectPageResponse](rpcResp.Data)
}
