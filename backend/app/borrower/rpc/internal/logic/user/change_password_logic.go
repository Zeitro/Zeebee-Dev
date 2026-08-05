package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *proto.ChangePasswordRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	input := dto.UserChangePasswordRequestDto{
		OldPassword: in.Data.OldPassword,
		NewPassword: in.Data.NewPassword,
		UserId:      claim.UserId,
		UserType:    claim.UserType,
	}
	err = l.baseLogic.UserLogic.ChangePassword(l.ctx, input)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
