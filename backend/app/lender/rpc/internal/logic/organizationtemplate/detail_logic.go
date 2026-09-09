package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
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

func (l *DetailLogic) Detail(in *proto.IdRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	detail, err := l.svcCtx.TenantClient.OrganizationTemplate.Detail(l.ctx, &tenant.OrganizationTemplateDetailRequest{Data: &tenant.OrganizationTemplateDetailRequestData{
		UserType:       string(constants.UserTypeLender),
		Id:             in.Data,
		OrganizationId: claim.OrganizationId,
	}})
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationTemplateDetailResponse{
		Data: detail.Data,
	}, nil
}
