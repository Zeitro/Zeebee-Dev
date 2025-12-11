package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type EditProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditProfileLogic {
	return &EditProfileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// rpc SendResetPasswordEmail(SendResetPasswordEmailRequest) returns(SendResetPasswordEmailResponse);
func (l *EditProfileLogic) EditProfile(in *admin.UserProfileEditRequest) (resp *admin.UserinfoResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
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
	data := &admin.UserinfoResponseData{}
	err = copier.CopyWithOption(data, token, copier.Option{
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	return &admin.UserinfoResponse{
		Data: data,
	}, nil
}
