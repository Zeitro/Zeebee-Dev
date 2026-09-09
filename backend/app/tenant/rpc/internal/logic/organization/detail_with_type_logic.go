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

type DetailWithTypeLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDetailWithTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailWithTypeLogic {
	return &DetailWithTypeLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// DetailWithType returns an organization by id, scoped to a platform type.
func (l *DetailWithTypeLogic) DetailWithType(in *tenant.OrganizationDetailWithTypeRequest) (resp *tenant.OrganizationDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	result, err := l.baseLogic.OrganizationLogic.DetailWithType(l.ctx, in.Data.Id, userType)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.OrganizationEntityDto, proto.OrganizationEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationDetailResponse{Data: data}, nil
}
