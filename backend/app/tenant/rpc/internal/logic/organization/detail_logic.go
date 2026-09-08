package organizationlogic

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

type DetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Detail returns an organization by id.
func (l *DetailLogic) Detail(in *tenant.OrganizationDetailRequest) (resp *tenant.OrganizationDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	result, err := l.baseLogic.OrganizationLogic.Detail(l.ctx, in.Data.Id)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.OrganizationEntityDto, proto.OrganizationEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationDetailResponse{Data: data}, nil
}
