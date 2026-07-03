package pricing_plan

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

type PricingPlanCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewPricingPlanCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanCreateLogic {
	return &PricingPlanCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanCreateLogic) PricingPlanCreate(req *types.PricingPlanCreateRequest) (resp *types.PricingPlanEntity, err error) {

	// Convert types to pb request
	rpcReq := &admin.PricingPlanCreateRequest{}
	rpcReq.Data, err = copierx.Trans[types.PricingPlanCreateRequest, admin.PricingPlanCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[proto.PricingPlanEntity, types.PricingPlanEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
