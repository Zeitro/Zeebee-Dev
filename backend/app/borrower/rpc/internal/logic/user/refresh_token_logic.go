package userlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

func (l *RefreshTokenLogic) RefreshToken(in *borrower.RefreshTokenRequest) (resp *borrower.RefreshTokenResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).Error(err)
		}
	}()
	var claims *utils.JWTClaims
	// parse refresh token
	claims, err = utils.ParseToken(in.RefreshToken, l.svcCtx.Config.JwtAuth.RefreshSecret)
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

	return &borrower.RefreshTokenResponse{
		AccessToken:       accessToken,
		AccessTokenExpire: l.svcCtx.Config.JwtAuth.AccessExpire,
	}, nil
}
