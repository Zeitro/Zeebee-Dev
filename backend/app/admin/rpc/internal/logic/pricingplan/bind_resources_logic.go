package pricingplanlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type BindResourcesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewBindResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindResourcesLogic {
	return &BindResourcesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BindResourcesLogic) BindResources(in *admin.PricingPlanBindResourcesRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	err = l.baseLogic.PricingPlanLogic.BindResources(l.ctx, in.Data.Id, in.Data.ResourceIds)
	return nil, errorx.WithStack(err)
}
