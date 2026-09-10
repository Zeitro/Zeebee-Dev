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

type AllLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// All returns every organization for a platform type.
func (l *AllLogic) All(in *tenant.OrganizationAllRequest) (resp *tenant.OrganizationAllResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	all, err := l.baseLogic.OrganizationLogic.All(l.ctx, userType)
	if err != nil {
		return
	}

	data, err := copierx.Trans[[]*dto.OrganizationEntityDto, []*proto.OrganizationEntity](&all, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationAllResponse{Data: *data}, nil
}
