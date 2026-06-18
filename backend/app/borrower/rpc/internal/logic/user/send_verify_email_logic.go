package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/limiter"
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

func (l *SendVerifyEmailLogic) SendVerifyEmail(in *borrower.SendVerifyEmailRequest) (resp *borrower.SendVerifyEmailResponse, err error) {
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
	key := stringx.Format(constants.RegisterVerifyEmail, constants.UserTypeBorrower, email)
	ok, err := global.LimiterMgr.Limit(l.ctx, limiter.StrategySendVerifyEmail, key)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !ok {
		return nil, errorx.Errorf("Email %s try too many times", email)
	}

	// check user status
	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, email, constants.UserTypeBorrower)
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
			err = errorx.Wrap(err, "Send register verify email error")
			l.WithContext(ctx).Error(err)
		}
	}()

	return &borrower.SendVerifyEmailResponse{}, nil
}
