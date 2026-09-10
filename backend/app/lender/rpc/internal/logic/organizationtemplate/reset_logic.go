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

type ResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResetLogic) Reset(in *proto.IdRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	organizationTemplate, err := l.svcCtx.TenantClient.OrganizationTemplate.Reset(l.ctx, &tenant.OrganizationTemplateResetRequest{Data: &tenant.OrganizationTemplateResetRequestData{
		UserType:       string(constants.UserTypeLender),
		Id:             in.Data,
		OrganizationId: claim.OrganizationId,
	}})
	if err != nil {
		return nil, err
	}
	return &lender.OrganizationTemplateDetailResponse{
		Data: organizationTemplate.Data,
	}, nil
}
