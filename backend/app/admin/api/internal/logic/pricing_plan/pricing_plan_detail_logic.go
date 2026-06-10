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

type PricingPlanDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewPricingPlanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanDetailLogic {
	return &PricingPlanDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanDetailLogic) PricingPlanDetail(req *types.PathIdRequest) (resp *types.PricingPlanDetailResponse, err error) {

	// Call RPC server
	rpcResp, err := l.svcCtx.PricingPlanRpcClient.Detail(l.ctx, &admin.PricingPlanDetailRequest{
		Data: &admin.PricingPlanDetailRequestData{Id: req.Id},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.PricingPlanEntity, types.PricingPlanDetailResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
