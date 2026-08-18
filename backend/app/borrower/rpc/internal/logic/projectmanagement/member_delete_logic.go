package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type MemberDeleteLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewMemberDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberDeleteLogic {
	return &MemberDeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *MemberDeleteLogic) MemberDelete(in *borrower.ProjectManagementMemberDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	if err = l.baseLogic.ProjectManagementLogic.RemoveProjectMember(l.ctx, actor, in.Data.ProjectId, in.Data.OrganizationUserId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
