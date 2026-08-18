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

type ProjectManagementTaskPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project tasks
func NewProjectManagementTaskPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskPageLogic {
	return &ProjectManagementTaskPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskPageLogic) ProjectManagementTaskPage(req *types.ProjectManagementTaskPageRequest) (resp *types.ProjectManagementTaskPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskPageRequest, borrower.ProjectManagementTaskPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.TaskPage(l.ctx, &borrower.ProjectManagementTaskPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskPageResponseData, types.ProjectManagementTaskPageResponse](rpcResp.Data)
}
