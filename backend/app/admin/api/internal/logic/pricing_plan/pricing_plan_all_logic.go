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

type PricingPlanAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewPricingPlanAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanAllLogic {
	return &PricingPlanAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanAllLogic) PricingPlanAll(req *types.PricingPLanAllrequest) (resp *types.PricingPlanAllResponse, err error) {
	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanRpcClient.All(l.ctx, &admin.PricingPlanAllRequest{
		Data: &admin.PricingPlanAllRequestData{Type: req.Type},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	list, err := copierx.Trans[[]*admin.PricingPlanEntity, []*types.PricingPlanEntity](&rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &types.PricingPlanAllResponse{List: *list}, nil
}
