package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *proto.PathIdRequest) (resp *admin.OrganizationSubscriptionDetailResponse, err error) {

	result, err := l.svcCtx.TenantClient.OrganizationSubscription.Detail(l.ctx, &tenant.OrganizationSubscriptionDetailRequest{Data: &tenant.OrganizationSubscriptionDetailRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data.Id,
	}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &admin.OrganizationSubscriptionDetailResponse{
		Data: result.Data,
	}, nil
}
