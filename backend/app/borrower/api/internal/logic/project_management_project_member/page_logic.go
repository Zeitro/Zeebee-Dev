package project_management_project_member

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

// page project members
func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(req *types.ProjectManagementMemberPageRequest) (resp *types.ProjectManagementMemberPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMemberPageRequest, borrower.ProjectManagementMemberPageRequestData](req)
	if err != nil {
		return nil, err
	}
	data.ProjectId = req.ProjectId
	rpcResp, err := l.svcCtx.ProjectManagementProjectMemberRpcClient.Page(l.ctx, &borrower.ProjectManagementMemberPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementMemberPageResponseData, types.ProjectManagementMemberPageResponse](rpcResp.Data)
}
