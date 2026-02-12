package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewOrganizationUserSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserSwitchStatusLogic {
	return &OrganizationUserSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserSwitchStatusLogic) OrganizationUserSwitchStatus(req *types.OrganizationUserSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.OrganizationUserSwitchStatusRequest, lender.OrganizationUserUpdateStatusRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcReq := &lender.OrganizationUserUpdateStatusRequest{
		Data: data,
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
