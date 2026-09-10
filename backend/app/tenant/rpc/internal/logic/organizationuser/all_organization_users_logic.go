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

type AllOrganizationUsersLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAllOrganizationUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllOrganizationUsersLogic {
	return &AllOrganizationUsersLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// AllOrganizationUsers returns every member of an organization.
func (l *AllOrganizationUsersLogic) AllOrganizationUsers(in *tenant.OrganizationUserAllOrganizationUsersRequest) (resp *tenant.OrganizationUserAllOrganizationUsersResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	all, err := l.baseLogic.OrganizationUserLogic.AllOrganizationUsers(l.ctx, in.Data.OrganizationId)
	if err != nil {
		return
	}

	data, err := copierx.Trans[[]*dto.OrganizationUserEntityDto, []*proto.OrganizationUserEntity](&all, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationUserAllOrganizationUsersResponse{Data: *data}, nil
}
