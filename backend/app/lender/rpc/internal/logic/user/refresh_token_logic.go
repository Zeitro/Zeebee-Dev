package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
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

func (l *RefreshTokenLogic) RefreshToken(in *lender.RefreshTokenRequest) (resp *lender.UserJwtResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	token, err := l.baseLogic.UserLogic.RefreshToken(l.ctx, in.Data.RefreshToken, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[bo.UserJwtTokenBo, lender.UserJwtResponseData](token)
	if err != nil {
		return nil, err
	}

	return &lender.UserJwtResponse{
		Data: data,
	}, nil
}
