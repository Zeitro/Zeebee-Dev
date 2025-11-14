package user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"

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
	rpcResp, err := l.svcCtx.UserRpcClient.RefreshToken(l.ctx, &lender.RefreshTokenRequest{
		Data: &lender.RefreshTokenRequestData{
			RefreshToken: req.RefreshToken,
		},
	})
	if err != nil {
		return nil, errorx.New("Refresh Token Error")
	}

	return &types.UserRefreshTokenResponse{
		AccessToken:        rpcResp.Data.AccessToken,
		AccessTokenExpire:  rpcResp.Data.AccessTokenExpire,
		RefreshToken:       rpcResp.Data.RefreshToken,
		RefreshTokenExpire: rpcResp.Data.RefreshTokenExpire,
	}, nil
}
