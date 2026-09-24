package projectmanagementprojectmemberlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type UpdateRoleLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *borrower.ProjectManagementMemberRoleUpdateRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.ProjectManagementMemberRoleUpdateRequestData, dto.ProjectManagementMemberDto](in.Data)
	if err != nil {
		return nil, err
	}
	if err = l.baseLogic.ProjectManagementMemberLogic.Update(l.ctx, claim, *input); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
