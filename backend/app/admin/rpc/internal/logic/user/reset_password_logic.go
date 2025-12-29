package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/errorx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// rpc VerifyEmail(VerifyEmailRequest) returns(VerifyEmailResponse);
func (l *ResetPasswordLogic) ResetPassword(in *admin.ResetPasswordRequest) (resp *admin.ResetPasswordResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
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
	data := &admin.UserLoginResponseData{}
	err = copier.CopyWithOption(&data, token, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.ResetPasswordResponse{
		Data: data,
	}, nil
}
