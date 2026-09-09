package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateExpireDateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateExpireDateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExpireDateLogic {
	return &UpdateExpireDateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateExpireDateLogic) UpdateExpireDate(in *admin.OrganizationSubscriptionUpdateExpireDateRequest) (resp *proto.EmptyResponse, err error) {

	input, err := copierx.Trans[admin.OrganizationSubscriptionUpdateExpireDateRequestData, dto.OrganizationSubscriptionUpdateExpireDateRequestDto](in.Data, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	err = l.baseLogic.OrganizationSubscriptionLogic.UpdateExpireDate(l.ctx, *input)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &proto.EmptyResponse{}, nil
}
