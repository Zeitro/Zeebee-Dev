package userlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	logic2 "zeebee-backend/domains/zeebee/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// rpc Register(UserRegisterRequest) returns(UserRegisterResponse);
func (l *RefreshTokenLogic) RefreshToken(in *admin.RefreshTokenRequest) (resp *admin.RefreshTokenResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var claims *utils.JWTClaims
	// parse refresh token
	claims, err = utils.ParseToken(in.Data.RefreshToken, l.svcCtx.Config.JwtAuth.RefreshSecret)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	id := claims.UserId

	user, err := l.svcCtx.UserRepo.GetUser(l.ctx, id, constants.UserTypeAdmin)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user.Status != constants.UserStatusActive {
		return nil, errorx.New("User Is Not Active")
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, logic2.JwtConfig{
		AccessSecret:  l.svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: l.svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
	}, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}
	return &admin.RefreshTokenResponse{
		Data: &admin.RefreshTokenResponseData{
			AccessToken:       token.AccessToken,
			AccessTokenExpire: token.AccessTokenExpire,
		},
	}, nil
}
