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

// EditProfile updates the caller's name and phone fields. user_id comes from
// the caller's JWT claim; organization_id is resolved from the user's default
// organization.
func (l *EditProfileLogic) EditProfile(in *tenant.EditProfileRequest) (resp *proto.UserinfoResponse, err error) {
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserRepo.GetUser(l.ctx, claim.UserId, claim.UserType, false)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New("user not found")
	}

	input := dto.OrganizationUserProfileEditRequestDto{
		OrganizationId: user.DefaultOrganizationId,
		UserId:         claim.UserId,
		UserType:       claim.UserType,
		FirstName:      in.Data.FirstName,
		MidName:        in.Data.MidName,
		LastName:       in.Data.LastName,
		Phone:          in.Data.Phone,
	}

	if err = l.baseLogic.OrganizationUserLogic.UpdateProfile(l.ctx, input); err != nil {
		return nil, err
	}

	info, err := l.baseLogic.UserLogic.Info(l.ctx, claim.UserId, user.DefaultOrganizationId, claim.UserType)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.UserinfoDto, proto.UserinfoResponseData](info)
	if err != nil {
		return nil, err
	}
	return &proto.UserinfoResponse{Data: data}, nil
}
