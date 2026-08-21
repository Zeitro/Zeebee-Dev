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

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// add project members
func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.ProjectManagementMemberAddRequest) (resp *types.EmptyResponse, err error) {
	members, err := copierx.Trans[[]*types.ProjectManagementMemberInput, []*borrower.ProjectManagementMemberInput](&req.Members)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProjectManagementProjectMemberRpcClient.Add(l.ctx, &borrower.ProjectManagementMemberAddRequest{Data: &borrower.ProjectManagementMemberAddRequestData{ProjectId: req.ProjectId, Members: *members}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
