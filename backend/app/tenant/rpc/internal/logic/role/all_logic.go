package rolelogic

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
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// All returns every role of the organization (role type scoped by user_type).
func (l *AllLogic) All(in *tenant.RoleAllRequest) (resp *tenant.RoleAllResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.RoleLogic.All(l.ctx, in.Data.OrganizationId, userType)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.RoleEntityDto, []*proto.RoleEntity](&result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.RoleAllResponse{Data: *data}, nil
}
