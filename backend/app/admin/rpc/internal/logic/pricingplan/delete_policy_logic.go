package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDeletePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePolicyLogic {
	return &DeletePolicyLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeletePolicyLogic) DeletePolicy(in *admin.PricingPlanPolicyDeleteRequest) (resp *admin.EmptyResponse, err error) {

	err = l.baseLogic.PricingPlanPolicyLogic.Delete(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &admin.EmptyResponse{}, nil
}
