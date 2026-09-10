package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateProfileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// UpdateProfile updates the requesting member's own profile fields.
func (l *UpdateProfileLogic) UpdateProfile(in *tenant.OrganizationUserUpdateProfileRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserProfileEditRequestDto{
		OrganizationId: in.Data.OrganizationId,
		UserId:         in.Data.UserId,
		UserType:       userType,
		FirstName:      in.Data.FirstName,
		MidName:        in.Data.MidName,
		LastName:       in.Data.LastName,
		Phone:          in.Data.Phone,
	}
	if err = l.baseLogic.OrganizationUserLogic.UpdateProfile(l.ctx, input); err != nil {
		return
	}
	return &proto.EmptyResponse{}, nil
}
