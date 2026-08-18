package project_management

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementMemberDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// remove project member
func NewProjectManagementMemberDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementMemberDeleteLogic {
	return &ProjectManagementMemberDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementMemberDeleteLogic) ProjectManagementMemberDelete(req *types.ProjectManagementMemberDeleteRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.ProjectManagementRpcClient.MemberDelete(l.ctx, &borrower.ProjectManagementMemberDeleteRequest{Data: &borrower.ProjectManagementMemberDeleteRequestData{ProjectId: req.Id, OrganizationUserId: req.UserId}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
