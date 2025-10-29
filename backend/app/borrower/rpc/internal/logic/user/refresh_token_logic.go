package userlogic

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

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
		return nil, errors.WithStack(err)
	}

	id, err := strconv.Atoi(claims.UserID)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserRepo.GetBorrowerUser(l.ctx, int64(id))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if user.Status != constants.UserStatusActive {
		return nil, errors.WithStack(errors.New("user is not active"))
	}

	accessToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.Borrower.Id)),
		0,
		user.Email,
		"",
		string(constants.UserTypeBorrower),
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &borrower.RefreshTokenResponse{
		Data: &borrower.RefreshTokenResponseData{
			AccessToken:       accessToken,
			AccessTokenExpire: l.svcCtx.Config.JwtAuth.AccessExpire,
		},
	}, nil
}
