package pricing_plan

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PricingPlanSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewPricingPlanSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanSwitchStatusLogic {
	return &PricingPlanSwitchStatusLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanSwitchStatusLogic) PricingPlanSwitchStatus(req *types.PricingPlanSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	// Convert types to pb request
	rpcReq := &admin.PricingPlanSwitchStatusRequest{}
	rpcReq.Data, err = copierx.Trans[types.PricingPlanSwitchStatusRequest, admin.PricingPlanSwitchStatusRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcReq.Data.Id = req.Id

	// Call RPC server
	_, err = l.svcCtx.PricingPlanRpcClient.SwitchStatus(l.ctx, rpcReq)
	return nil, errorx.WithStack(err)
}
