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

type UpdateOrganizationUserLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateOrganizationUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrganizationUserLogic {
	return &UpdateOrganizationUserLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// UpdateOrganizationUser updates a member's profile and roles.
func (l *UpdateOrganizationUserLogic) UpdateOrganizationUser(in *tenant.OrganizationUserUpdateRequest) (resp *tenant.OrganizationUserDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserUpdateRequestDto{
		Id:             in.Data.Id,
		Phone:          in.Data.Phone,
		FirstName:      in.Data.FirstName,
		MidName:        in.Data.MidName,
		LastName:       in.Data.LastName,
		Remark:         in.Data.Remark,
		RoleIds:        in.Data.RoleIds,
		UserType:       userType,
		OrganizationId: in.Data.OrganizationId,
	}
	result, err := l.baseLogic.OrganizationUserLogic.UpdateOrganizationUser(l.ctx, input)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.OrganizationUserEntityDto, proto.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationUserDetailResponse{Data: data}, nil
}
