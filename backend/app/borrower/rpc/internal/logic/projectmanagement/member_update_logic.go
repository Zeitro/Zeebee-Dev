package projectmanagementlogic

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

type MemberUpdateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewMemberUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberUpdateLogic {
	return &MemberUpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *MemberUpdateLogic) MemberUpdate(in *borrower.ProjectManagementMemberUpdateRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	if in.Data.TransferOwnerToUserId != 0 {
		err = l.baseLogic.ProjectManagementLogic.TransferProjectOwner(l.ctx, actor, in.Data.ProjectId, in.Data.TransferOwnerToUserId)
	} else {
		input, convertErr := copierx.Trans[borrower.ProjectManagementMemberUpdateRequestData, dto.ProjectManagementMemberDto](in.Data)
		if convertErr != nil {
			return nil, convertErr
		}
		err = l.baseLogic.ProjectManagementLogic.UpdateProjectMemberRole(l.ctx, actor, *input)
	}
	if err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
