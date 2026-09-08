package project_management_project_member

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// remove project member
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.ProjectManagementMemberDeleteRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.ProjectManagementProjectMemberRpcClient.Delete(l.ctx, &borrower.ProjectManagementMemberDeleteRequest{Data: &borrower.ProjectManagementMemberDeleteRequestData{ProjectId: req.ProjectId, MemberId: req.Id}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
