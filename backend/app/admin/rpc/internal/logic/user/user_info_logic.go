package userlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
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

func (l *UserInfoLogic) UserInfo(in *proto.EmptyRequest) (resp *admin.UserinfoResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	user, err := l.baseLogic.UserLogic.Info(l.ctx, claim.UserId, claim.OrganizationId, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserinfoDto, admin.UserinfoResponseData](user)
	if err != nil {
		return nil, err
	}
	userInfo := &admin.UserinfoResponse{
		Data: data,
	}
	return userInfo, nil
}
