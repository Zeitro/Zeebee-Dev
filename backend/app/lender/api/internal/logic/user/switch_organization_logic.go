// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchOrganizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch lender company
func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(req *types.UserSwitchOrganizationRequest) (resp *types.UserRefreshTokenResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpcClient.SwitchOrganization(l.ctx, &lender.SwitchOrganizationRequest{
		Data: &lender.SwitchOrganizationRequestData{
			OrganizationId: req.OrganizationId,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.UserRefreshTokenResponse{
		AccessToken:        rpcResp.Data.AccessToken,
		AccessTokenExpire:  rpcResp.Data.AccessTokenExpire,
		RefreshToken:       rpcResp.Data.RefreshToken,
		RefreshTokenExpire: rpcResp.Data.RefreshTokenExpire,
	}, nil
}
