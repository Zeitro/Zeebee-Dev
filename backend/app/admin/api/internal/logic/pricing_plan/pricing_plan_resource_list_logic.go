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

type PricingPlanResourceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// resource list
func NewPricingPlanResourceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanResourceListLogic {
	return &PricingPlanResourceListLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanResourceListLogic) PricingPlanResourceList(req *types.PathIdRequest) (resp *types.PricingPlanResourceListResponse, err error) {

	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanRpcClient.ResourceList(l.ctx, &admin.PricingPlanResourceListRequest{
		Data: &admin.PricingPlanResourceListRequestData{Id: req.Id},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.PricingPlanResourceListResponse, types.PricingPlanResourceListResponse](rpcResp)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
