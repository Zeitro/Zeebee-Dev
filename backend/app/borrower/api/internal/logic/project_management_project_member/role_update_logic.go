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

type RoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update role
func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleUpdateLogic) RoleUpdate(req *types.ProjectManagementMemberRoleUpdateRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMemberRoleUpdateRequest, borrower.ProjectManagementMemberRoleUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	data.MemberId = req.Id
	if _, err = l.svcCtx.ProjectManagementProjectMemberRpcClient.UpdateRole(l.ctx, &borrower.ProjectManagementMemberRoleUpdateRequest{Data: data}); err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
