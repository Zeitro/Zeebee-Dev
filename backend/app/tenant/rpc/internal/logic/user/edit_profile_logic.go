package userlogic

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

type EditProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditProfileLogic {
	return &EditProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// EditProfile updates the user's name fields. Phone is preserved (the tenant
// profile contract is name-only); organization_id is resolved from the user's
// default organization.
func (l *EditProfileLogic) EditProfile(in *tenant.EditProfileRequest) (resp *proto.UserinfoResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserRepo.GetUser(l.ctx, in.Data.UserId, userType, false)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New("user not found")
	}
	orgUser, err := l.svcCtx.OrganizationUserRepo.GetOrganizationUserByUserId(l.ctx, user.DefaultOrganizationId, in.Data.UserId)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.EditProfileRequestData, dto.OrganizationUserProfileEditRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.OrganizationId = user.DefaultOrganizationId
	input.UserId = in.Data.UserId
	input.UserType = userType
	input.Phone = orgUser.Phone // preserve — tenant profile edit is name-only

	if err = l.baseLogic.OrganizationUserLogic.UpdateProfile(l.ctx, *input); err != nil {
		return nil, err
	}

	info, err := l.baseLogic.UserLogic.Info(l.ctx, in.Data.UserId, user.DefaultOrganizationId, userType)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.UserinfoDto, proto.UserinfoResponseData](info)
	if err != nil {
		return nil, err
	}
	return &proto.UserinfoResponse{Data: data}, nil
}
