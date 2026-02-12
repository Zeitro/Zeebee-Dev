package organization_setting

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationSettingSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewOrganizationSettingSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSettingSwitchStatusLogic {
	return &OrganizationSettingSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSettingSwitchStatusLogic) OrganizationSettingSwitchStatus(req *types.OrganizationSettingSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &lender.OrganizationSettingSwitchStatusRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationSettingSwitchStatusRequest, lender.OrganizationSettingSwitchStatusRequestData](req)
	if err != nil {
		return
	}

	_, err = l.svcCtx.OrganizationSettingRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return
	}

	return &types.EmptyResponse{}, nil
}
