package user

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// refresh
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.UserRefreshTokenRequest) (resp *types.UserRefreshTokenResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpcClient.RefreshToken(l.ctx, &borrower.RefreshTokenRequest{
		Data: &borrower.RefreshTokenRequestData{
			RefreshToken: req.RefreshToken,
		},
	})
	if err != nil {
		return nil, errors.WithStack(errors.New("refresh token error"))
	}

	return &types.UserRefreshTokenResponse{
		AccessToken:       rpcResp.Data.AccessToken,
		AccessTokenExpire: rpcResp.Data.AccessTokenExpire,
	}, nil
}