package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/errorx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// rpc VerifyEmail(VerifyEmailRequest) returns(VerifyEmailResponse);
func (l *ResetPasswordLogic) ResetPassword(in *admin.ResetPasswordRequest) (resp *admin.ResetPasswordResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	input := dto.UserResetPasswordByEmailRequestDto{
		UserType:    constants.UserTypeAdmin,
		NewPassword: in.Data.Password,
		Email:       in.Data.Email,
		Code:        in.Data.Code,
	}

	token, err := l.baseLogic.UserLogic.ResetPasswordWithEmail(l.ctx, input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserJwtTokenDto, admin.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}

	return &admin.ResetPasswordResponse{
		Data: data,
	}, nil
}
