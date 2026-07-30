package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateOrderLogic) UpdateOrder(in *admin.OrganizationUpdateOrderRequest) (resp *proto.EmptyResponse, err error) {

	err = l.baseLogic.OrganizationLogic.UpdateOrder(l.ctx, in.Data.Id, int(in.Data.Order))
	if err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
