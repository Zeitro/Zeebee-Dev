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

type ProjectManagementTaskTypePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project task types
func NewProjectManagementTaskTypePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskTypePageLogic {
	return &ProjectManagementTaskTypePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskTypePageLogic) ProjectManagementTaskTypePage(req *types.ProjectManagementTaskTypePageRequest) (resp *types.ProjectManagementTaskTypePageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTypePageRequest, borrower.ProjectManagementTaskTypePageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskTypePage(l.ctx, &borrower.ProjectManagementTaskTypePageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskTypePageResponseData, types.ProjectManagementTaskTypePageResponse](rpcResp.Data)
}
