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

type ProjectManagementMemberUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update project member
func NewProjectManagementMemberUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementMemberUpdateLogic {
	return &ProjectManagementMemberUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementMemberUpdateLogic) ProjectManagementMemberUpdate(req *types.ProjectManagementMemberUpdateRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMemberUpdateRequest, borrower.ProjectManagementMemberUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId, data.OrganizationUserId = req.Id, req.UserId
	_, err = l.svcCtx.ProjectManagementRpcClient.MemberUpdate(l.ctx, &borrower.ProjectManagementMemberUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
