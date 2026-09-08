package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type AllForOptionLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAllForOptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllForOptionLogic {
	return &AllForOptionLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// AllForOption returns lightweight organization options filtered by status.
func (l *AllForOptionLogic) AllForOption(in *tenant.OrganizationAllForOptionRequest) (resp *tenant.OrganizationOptionsResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	all, err := l.baseLogic.OrganizationLogic.AllForOption(l.ctx, userType, constants.OrganizationStatus(in.Data.Status))
	if err != nil {
		return
	}

	data, err := copierx.Trans[[]*dto.OrganizationOptionDto, []*proto.OrganizationOptionsEntity](&all)
	if err != nil {
		return
	}
	return &tenant.OrganizationOptionsResponse{Data: *data}, nil
}
