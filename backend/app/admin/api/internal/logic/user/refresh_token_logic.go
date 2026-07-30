package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// refresh
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.UserRefreshTokenRequest) (resp *types.UserRefreshTokenResponse, err error) {
	resp = &types.UserRefreshTokenResponse{}
	rpcReq := &admin.RefreshTokenRequest{
		Data: &admin.RefreshTokenRequestData{
			RefreshToken: req.RefreshToken,
		},
	}

	// call rpc
	rpcResp, err := l.svcCtx.UserRpcClient.RefreshToken(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[admin.RefreshTokenResponseData, types.UserRefreshTokenResponse](rpcResp.Data)
	return resp, err
}
