package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
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

func (l *EditProfileLogic) EditProfile(in *borrower.UserProfileEditRequest) (resp *proto.UserinfoResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
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
	token, err := l.baseLogic.UserLogic.Info(l.ctx, claim.UserId, claim.OrganizationUserId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &proto.UserinfoResponseData{}
	data, err = copierx.Trans[dto.UserinfoDto, proto.UserinfoResponseData](token)
	if err != nil {
		return nil, err
	}

	return &proto.UserinfoResponse{
		Data: data,
	}, nil
}
