package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// RefreshToken re-issues a TokenSet from a refresh token. admin returns the full
// TokenSet (ODD-6); the platform key comes from user_type.
func (l *RefreshTokenLogic) RefreshToken(in *tenant.RefreshTokenRequest) (resp *proto.UserLoginResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	token, err := l.baseLogic.UserLogic.RefreshToken(l.ctx, in.Data.RefreshToken, userType)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}
	return &proto.UserLoginResponse{Data: data}, nil
}
