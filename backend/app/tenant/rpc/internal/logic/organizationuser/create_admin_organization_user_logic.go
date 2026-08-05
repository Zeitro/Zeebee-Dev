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

type CreateAdminOrganizationUserLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewCreateAdminOrganizationUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminOrganizationUserLogic {
	return &CreateAdminOrganizationUserLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// CreateAdminOrganizationUser creates an admin-provisioned member with credentials.
func (l *CreateAdminOrganizationUserLogic) CreateAdminOrganizationUser(in *tenant.OrganizationUserCreateAdminRequest) (resp *tenant.OrganizationUserDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserCreateAdminRequestDto{
		Email:          in.Data.Email,
		Phone:          in.Data.Phone,
		FirstName:      in.Data.FirstName,
		LastName:       in.Data.LastName,
		Password:       in.Data.Password,
		Remark:         in.Data.Remark,
		RoleIds:        in.Data.RoleIds,
		UserType:       userType,
		OrganizationId: in.Data.OrganizationId,
	}
	result, err := l.baseLogic.OrganizationUserLogic.CreateAdminOrganizationUser(l.ctx, input)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.OrganizationUserEntityDto, proto.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationUserDetailResponse{Data: data}, nil
}
