// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

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

func (l *SwitchOrganizationLogic) SwitchOrganization(req *types.UserSwitchOrganizationRequest) (resp *types.UserJwtResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpcClient.SwitchOrganization(l.ctx, &lender.SwitchOrganizationRequest{
		Data: &lender.SwitchOrganizationRequestData{
			OrganizationId: req.OrganizationId,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err = copierx.Trans[lender.UserJwtResponseData, types.UserJwtResponse](rpcResp.Data)

	return resp, err
}
