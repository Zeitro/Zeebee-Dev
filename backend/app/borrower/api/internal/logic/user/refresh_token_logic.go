package user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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

func (l *RefreshTokenLogic) RefreshToken(req *types.UserRefreshTokenRequest) (resp *types.UserLoginResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpcClient.RefreshToken(l.ctx, &proto.RefreshTokenRequest{
		Data: &proto.RefreshTokenRequestData{
			RefreshToken: req.RefreshToken,
		},
	})
	if err != nil {
		return nil, errorx.New("refresh token error")
	}

	return &types.UserLoginResponse{
		AccessToken:        rpcResp.Data.AccessToken,
		AccessTokenExpire:  rpcResp.Data.AccessTokenExpire,
		RefreshToken:       rpcResp.Data.RefreshToken,
		RefreshTokenExpire: rpcResp.Data.RefreshTokenExpire,
	}, nil
}
