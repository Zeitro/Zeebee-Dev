package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
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

func (l *RefreshTokenLogic) RefreshToken(in *borrower.RefreshTokenRequest) (resp *borrower.RefreshTokenResponse, err error) {
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

	user, err := l.svcCtx.UserRepo.GetBorrowerUser(l.ctx, int64(id))
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user.Status != constants.UserStatusActive {
		return nil, errorx.New("User Is Not Active")
	}

	userClaim := utils.JwtUser{
		UserId:         user.Id,
		OrganizationId: user.DefaultOrganizationId,
		Username:       user.FullName,
		Role:           nil,
		UserType:       constants.UserTypeBorrower,
	}

	accessToken, err := utils.GenerateToken(
		userClaim,
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &borrower.RefreshTokenResponse{
		Data: &borrower.RefreshTokenResponseData{
			AccessToken:       accessToken,
			AccessTokenExpire: l.svcCtx.Config.JwtAuth.AccessExpire,
		},
	}, nil
}
