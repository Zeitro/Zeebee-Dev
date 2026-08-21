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

type AddLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AddLogic) Add(in *borrower.ProjectManagementMemberAddRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := claimValue
	members, err := copierx.Trans[[]*borrower.ProjectManagementMemberInput, []dto.ProjectManagementMemberDto](&in.Data.Members)
	if err != nil {
		return nil, err
	}
	for index := range *members {
		(*members)[index].ProjectId = in.Data.ProjectId
	}
	if err = l.baseLogic.ProjectManagementMemberLogic.Add(l.ctx, claim, *members); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
