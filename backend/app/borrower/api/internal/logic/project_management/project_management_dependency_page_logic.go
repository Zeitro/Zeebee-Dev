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

type ProjectManagementDependencyPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page task dependencies
func NewProjectManagementDependencyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementDependencyPageLogic {
	return &ProjectManagementDependencyPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementDependencyPageLogic) ProjectManagementDependencyPage(req *types.ProjectManagementDependencyPageRequest) (resp *types.ProjectManagementDependencyPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementDependencyPageRequest, borrower.ProjectManagementDependencyPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.DependencyPage(l.ctx, &borrower.ProjectManagementDependencyPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementDependencyPageResponseData, types.ProjectManagementDependencyPageResponse](rpcResp.Data)
}
