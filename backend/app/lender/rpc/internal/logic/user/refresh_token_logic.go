package userlogic

import (
	"context"
	"strconv"

	"zeebee-backend/app/lender/rpc/internal/logic"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
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
		return nil, errorx.WithStack(err)
	}
	id, _ := strconv.ParseInt(claims.UserID, 10, 64)

	user, err := l.svcCtx.UserRepo.GetLenderUser(l.ctx, id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	if user.Status == constants.UserStatusInactive {
		return nil, errorx.New("User Is Inactive")
	}

	employee, err := l.baseLogic.CheckAndUpdateActiveLenderCompanyByEmployee(user.LenderEmployee)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	user.LenderEmployee = employee

	token, err := generateLenderToken(user, l.svcCtx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &lender.RefreshTokenResponse{
		Data: &lender.RefreshTokenResponseData{
			AccessToken:        token.AccessToken,
			AccessTokenExpire:  token.AccessTokenExpire,
			RefreshToken:       token.RefreshToken,
			RefreshTokenExpire: token.RefreshTokenExpire,
		},
	}, nil
}
