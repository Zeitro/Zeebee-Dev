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

type PricingPlanUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewPricingPlanUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanUpdateLogic {
	return &PricingPlanUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanUpdateLogic) PricingPlanUpdate(req *types.PricingPlanUpdateRequest) (resp *types.PricingPlanEntity, err error) {

	// Convert types to pb request
	rpcReq := &admin.PricingPlanUpdateRequest{}
	rpcReq.Data, err = copierx.Trans[types.PricingPlanUpdateRequest, admin.PricingPlanUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcReq.Data.Id = req.Id

	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.PricingPlanEntity, types.PricingPlanEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
