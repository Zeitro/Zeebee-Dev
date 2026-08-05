package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UserInfoLogic) UserInfo(in *proto.EmptyRequest) (resp *proto.UserinfoResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	user, err := l.baseLogic.UserLogic.Info(l.ctx, claim.UserId, claim.OrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &proto.UserinfoResponseData{}
	data, err = copierx.Trans[dto.UserinfoDto, proto.UserinfoResponseData](user)
	if err != nil {
		return nil, err
	}
	userInfo := &proto.UserinfoResponse{
		Data: data,
	}
	return userInfo, nil
}
