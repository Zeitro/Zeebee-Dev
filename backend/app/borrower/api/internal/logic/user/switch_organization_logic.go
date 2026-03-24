package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type SwitchOrganizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch organization
func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(req *types.UserSwitchOrganizationRequest) (resp *types.UserLoginResponse, err error) {
	resp = &types.UserLoginResponse{}
	rpcReq := &borrower.SwitchOrganizationRequest{
		Data: &borrower.SwitchOrganizationRequestData{},
	}
	rpcReq.Data.OrganizationId = req.OrganizationId

	rpcResp, err := l.svcCtx.UserRpcClient.SwitchOrganization(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[borrower.UserLoginResponseData, types.UserLoginResponse](rpcResp.Data)

	if err != nil {
		return nil, err
	}
	return resp, err
}
