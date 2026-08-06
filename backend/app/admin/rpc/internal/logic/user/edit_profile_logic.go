package userlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

// rpc SendResetPasswordEmail(SendResetPasswordEmailRequest) returns(SendResetPasswordEmailResponse);
func (l *EditProfileLogic) EditProfile(in *admin.UserProfileEditRequest) (resp *admin.UserinfoResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input := dto.OrganizationUserProfileEditRequestDto{
		OrganizationId: claim.OrganizationId,
		UserId:         claim.UserId,
		UserType:       claim.UserType,
		FirstName:      in.Data.FirstName,
		MidName:        "",
		LastName:       in.Data.LastName,
		Phone:          in.Data.Phone,
	}
	err = l.baseLogic.OrganizationUserLogic.UpdateProfile(l.ctx, input)
	if err != nil {
		return nil, err
	}
	token, err := l.baseLogic.UserLogic.Info(l.ctx, claim.UserId, claim.OrganizationId, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserinfoDto, admin.UserinfoResponseData](token)
	if err != nil {
		return nil, err
	}

	return &admin.UserinfoResponse{
		Data: data,
	}, nil
}
