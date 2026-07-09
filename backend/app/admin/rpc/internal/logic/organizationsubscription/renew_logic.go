package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type RenewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewLogic {
	return &RenewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *RenewLogic) Renew(in *admin.OrganizationSubscriptionRenewRequest) (resp *proto.EmptyResponse, err error) {
	result, err := l.baseLogic.OrganizationSubscriptionLogic.Renew(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	_, err = copierx.Trans[dto.OrganizationSubscriptionEntity, proto.OrganizationSubscriptionEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
