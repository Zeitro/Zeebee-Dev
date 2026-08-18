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

type ProjectManagementMemberPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project members
func NewProjectManagementMemberPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementMemberPageLogic {
	return &ProjectManagementMemberPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementMemberPageLogic) ProjectManagementMemberPage(req *types.ProjectManagementMemberPageRequest) (resp *types.ProjectManagementMemberPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMemberPageRequest, borrower.ProjectManagementMemberPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.MemberPage(l.ctx, &borrower.ProjectManagementMemberPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementMemberPageResponseData, types.ProjectManagementMemberPageResponse](rpcResp.Data)
}
