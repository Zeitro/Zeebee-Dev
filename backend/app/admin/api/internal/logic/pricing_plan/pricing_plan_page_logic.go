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

type PricingPlanPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewPricingPlanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanPageLogic {
	return &PricingPlanPageLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanPageLogic) PricingPlanPage(req *types.PricingPlanPageRequest) (resp *types.PricingPlanPageResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.PricingPlanPageRequest{}
	rpcReq.Data, err = copierx.Trans[types.PricingPlanPageRequest, admin.PricingPlanPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.PricingPlanPageResponseData, types.PricingPlanPageResponse](rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
