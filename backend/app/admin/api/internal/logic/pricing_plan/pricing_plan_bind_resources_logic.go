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

type PricingPlanBindResourcesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// bind resources
func NewPricingPlanBindResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanBindResourcesLogic {
	return &PricingPlanBindResourcesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanBindResourcesLogic) PricingPlanBindResources(req *types.PricingPlanBindResourcesRequest) (resp *types.PricingPlanResourceListResponse, err error) {
	// Convert types to pb request
	rpcReq := &admin.PricingPlanBindResourcesRequest{}
	rpcReq.Data, err = copierx.Trans[types.PricingPlanBindResourcesRequest, admin.PricingPlanBindResourcesRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcReq.Data.Id = req.Id

	// Call RPC server
	_, err = l.svcCtx.PricingPlanRpcClient.BindResources(l.ctx, rpcReq)
	return nil, errorx.WithStack(err)

}
