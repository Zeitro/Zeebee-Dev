package pricing_plan

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PricingPlanDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewPricingPlanDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PricingPlanDeleteLogic {
	return &PricingPlanDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PricingPlanDeleteLogic) PricingPlanDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {

	// Call RPC server
	_, err = l.svcCtx.PricingPlanRpcClient.Delete(l.ctx, &admin.PathIdRequest{
		Data: &admin.PathIdRequestData{Id: req.Id},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}
