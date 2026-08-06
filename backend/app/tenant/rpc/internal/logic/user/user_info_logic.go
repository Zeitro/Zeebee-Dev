package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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

// UserInfo returns the user profile for the active organization. organization_id
// is resolved from the user's default organization; the user type comes from the
// caller's JWT claim.
func (l *UserInfoLogic) UserInfo(in *tenant.UserInfoRequest) (resp *proto.UserinfoResponse, err error) {
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserRepo.GetUser(l.ctx, in.Data.UserId, claim.UserType, false)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New("user not found")
	}

	info, err := l.baseLogic.UserLogic.Info(l.ctx, in.Data.UserId, user.DefaultOrganizationId, claim.UserType)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserinfoDto, proto.UserinfoResponseData](info)
	if err != nil {
		return nil, err
	}
	return &proto.UserinfoResponse{Data: data}, nil
}
