package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
)

type DetailByPurposeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
	// repo      *repo.AllRepos
}

func NewDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByPurposeLogic {
	return &DetailByPurposeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailByPurposeLogic) DetailByPurpose(in *lender.OrganizationTemplateDetailByPurposeRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {

	claims, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	detail, err := l.svcCtx.TenantClient.OrganizationTemplate.DetailByPurpose(l.ctx, &tenant.OrganizationTemplateDetailByPurposeRequest{Data: &tenant.OrganizationTemplateDetailByPurposeRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: claims.OrganizationId,
		Purpose:        in.Data.Purpose,
	}})
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationTemplateDetailResponse{
		Data: detail.Data,
	}, nil
}
