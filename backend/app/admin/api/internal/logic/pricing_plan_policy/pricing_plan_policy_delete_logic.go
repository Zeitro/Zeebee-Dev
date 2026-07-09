package pricing_plan_policy

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type PricingPlanPolicyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewPricingPlanPolicyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanPolicyDeleteLogic {
	return &PricingPlanPolicyDeleteLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanPolicyDeleteLogic) PricingPlanPolicyDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {

	// Call RPC server
	_, err = l.svcCtx.PricingPlanPolicyRpcClient.Delete(l.ctx, &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{Id: req.Id},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}
