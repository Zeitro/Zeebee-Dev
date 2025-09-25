package userlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *lender.RefreshTokenRequest) (resp *lender.RefreshTokenResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var claims *utils.JWTClaims
	// parse refresh token
	claims, err = utils.ParseToken(in.Data.RefreshToken, l.svcCtx.Config.JwtAuth.RefreshSecret)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	accessToken, err := utils.GenerateToken(
		claims.UserID,
		claims.AccountId,
		claims.Username,
		claims.Role,
		claims.UserType,
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &lender.RefreshTokenResponse{
		Data: &lender.RefreshTokenResponseData{
			AccessToken:       accessToken,
			AccessTokenExpire: l.svcCtx.Config.JwtAuth.AccessExpire,
		},
	}, nil
}
