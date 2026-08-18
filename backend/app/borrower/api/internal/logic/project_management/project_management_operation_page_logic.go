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

type ProjectManagementOperationPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project operation timeline
func NewProjectManagementOperationPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementOperationPageLogic {
	return &ProjectManagementOperationPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementOperationPageLogic) ProjectManagementOperationPage(req *types.ProjectManagementOperationPageRequest) (resp *types.ProjectManagementOperationPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementOperationPageRequest, borrower.ProjectManagementOperationPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.OperationPage(l.ctx, &borrower.ProjectManagementOperationPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementOperationPageResponseData, types.ProjectManagementOperationPageResponse](rpcResp.Data)
}
