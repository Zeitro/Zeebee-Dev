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

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchStatusLogic) SwitchStatus(in *lender.OrganizationSettingSwitchStatusRequest) (resp *lender.OrganizationSettingDetailResponse, err error) {

	if in.Data == nil {
		return nil, errorx.New("Data is required")
	}

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.TenantClient.OrganizationSetting.SwitchStatus(l.ctx, &tenant.OrganizationSettingSwitchStatusRequest{Data: &tenant.OrganizationSettingSwitchStatusRequestData{
		UserType:       string(constants.UserTypeLender),
		Id:             in.Data.Id,
		OrganizationId: claim.OrganizationId,
		Enabled:        in.Data.Enabled,
	}})
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationSettingDetailResponse{}, nil
}
