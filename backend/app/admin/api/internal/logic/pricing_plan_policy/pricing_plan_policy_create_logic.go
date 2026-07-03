package pricing_plan_policy

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type PricingPlanPolicyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewPricingPlanPolicyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanPolicyCreateLogic {
	return &PricingPlanPolicyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanPolicyCreateLogic) PricingPlanPolicyCreate(req *types.PricingPlanPolicyCreateRequest) (resp *types.PricingPlanPolicyEntity, err error) {

	// Convert types to pb request
	rpcReq := &admin.PricingPlanPolicyCreateRequest{}
	rpcReq.Data, err = copierx.Trans[types.PricingPlanPolicyCreateRequest, admin.PricingPlanPolicyCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanPolicyRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[proto.PricingPlanPolicyEntity, types.PricingPlanPolicyEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
