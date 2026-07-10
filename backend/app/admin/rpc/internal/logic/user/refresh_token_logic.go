package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
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

// rpc Register(UserRegisterRequest) returns(UserRegisterResponse);
func (l *RefreshTokenLogic) RefreshToken(in *admin.RefreshTokenRequest) (resp *admin.RefreshTokenResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	token, err := l.baseLogic.UserLogic.RefreshToken(l.ctx, in.Data.RefreshToken, constants.UserTypeAdmin)
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
