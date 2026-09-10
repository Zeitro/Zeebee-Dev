package organizationsettinglogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailByPurposeLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByPurposeLogic {
	return &DetailByPurposeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// DetailByPurpose returns the organization setting for a given purpose.
func (l *DetailByPurposeLogic) DetailByPurpose(in *tenant.OrganizationSettingDetailByPurposeRequest) (resp *tenant.OrganizationSettingDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	result, err := l.baseLogic.OrganizationSettingLogic.DetailByPurpose(l.ctx, in.Data.OrganizationId, in.Data.Purpose)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.OrganizationSettingEntityDto, proto.OrganizationSettingEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationSettingDetailResponse{Data: data}, nil
}
