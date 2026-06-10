package pricing_plan_policy

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PricingPlanPolicyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewPricingPlanPolicyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanPolicyUpdateLogic {
	return &PricingPlanPolicyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanPolicyUpdateLogic) PricingPlanPolicyUpdate(req *types.PricingPlanPolicyUpdateRequest) (resp *types.PricingPlanPolicyEntityResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.PricingPlanPolicyUpdateRequest{}
	rpcReq.Data, err = copierx.Trans[types.PricingPlanPolicyUpdateRequest, admin.PricingPlanPolicyUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcReq.Data.Id = req.Id

	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanRpcClient.UpdatePolicy(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.PricingPlanPolicyEntity, types.PricingPlanPolicyEntityResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
