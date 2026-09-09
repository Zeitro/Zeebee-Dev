package organizationsettinglogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"

	"zeebee-backend/common/errorx"
)

type DetailByPurposeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByPurposeLogic {
	return &DetailByPurposeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailByPurposeLogic) DetailByPurpose(in *lender.OrganizationSettingDetailByPurposeRequest) (resp *lender.OrganizationSettingDetailResponse, err error) {

	if in.Data == nil {
		return nil, errorx.New("data is required")
	}

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.TenantClient.OrganizationSetting.DetailByPurpose(l.ctx, &tenant.OrganizationSettingDetailByPurposeRequest{Data: &tenant.OrganizationSettingDetailByPurposeRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: claim.OrganizationId,
		Purpose:        in.Data.Purpose,
	}})
	if err != nil {
		return nil, err
	}
	return &lender.OrganizationSettingDetailResponse{
		Data: result.Data,
	}, nil
}
