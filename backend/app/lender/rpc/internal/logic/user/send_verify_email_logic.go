package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/limiter"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendVerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyEmailLogic {
	return &SendVerifyEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendVerifyEmailLogic) SendVerifyEmail(in *lender.SendVerifyEmailRequest) (resp *lender.SendVerifyEmailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	email := in.Data.Email
	// check send email limit
	key := stringx.Format(constants.RegisterVerifyEmail, constants.UserTypeLender, email)
	ok, err := global.LimiterMgr.Limit(l.ctx, limiter.StrategySendVerifyEmail, key)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !ok {
		return nil, errorx.Errorf("Email %s Try Too Many Times", email)
	}

	// check user status
	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, email, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	if user.Status != constants.UserStatusUnverified {
		return nil, errorx.WithStack(errorx.New("User is verified"))
	}

	// send email
	go func() {
		// send account validate email
		ctx := context.WithoutCancel(l.ctx)
		err := sendRegisterVerifyUrlEmail(ctx, email, user.FirstName)
		if err != nil {
			l.WithContext(ctx).Error(err)
		}
		if err != nil {
			err = errorx.Wrap(err, "Send Register Verify Email Error")
			l.WithContext(ctx).Error(err)
		}
	}()

	return &lender.SendVerifyEmailResponse{}, nil
}
