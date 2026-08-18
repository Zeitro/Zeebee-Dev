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

type ProjectManagementMyTaskPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page my project tasks
func NewProjectManagementMyTaskPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementMyTaskPageLogic {
	return &ProjectManagementMyTaskPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementMyTaskPageLogic) ProjectManagementMyTaskPage(req *types.ProjectManagementMyTaskPageRequest) (resp *types.ProjectManagementTaskPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMyTaskPageRequest, borrower.ProjectManagementMyTaskPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.MyTaskPage(l.ctx, &borrower.ProjectManagementMyTaskPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementTaskPageResponseData, types.ProjectManagementTaskPageResponse](rpcResp.Data)
}
