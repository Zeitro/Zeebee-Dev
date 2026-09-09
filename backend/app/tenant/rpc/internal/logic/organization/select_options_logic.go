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

type SelectOptionsLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSelectOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectOptionsLogic {
	return &SelectOptionsLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// SelectOptions returns lightweight organization options for a platform type.
func (l *SelectOptionsLogic) SelectOptions(in *tenant.OrganizationSelectOptionsRequest) (resp *tenant.OrganizationOptionsResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	all, err := l.baseLogic.OrganizationLogic.SelectOptions(l.ctx, userType)
	if err != nil {
		return
	}

	data, err := copierx.Trans[[]*dto.OrganizationOptionDto, []*proto.OrganizationOptionsEntity](&all)
	if err != nil {
		return
	}
	return &tenant.OrganizationOptionsResponse{Data: *data}, nil
}
