package organizationuserlogic

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

// Detail returns an organization member by id (org-scoped).
func (l *DetailLogic) Detail(in *tenant.OrganizationUserDetailRequest) (resp *tenant.OrganizationUserDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	result, err := l.baseLogic.OrganizationUserLogic.Detail(l.ctx, in.Data.Id, in.Data.OrganizationId)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.OrganizationUserEntityDto, proto.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationUserDetailResponse{Data: data}, nil
}
