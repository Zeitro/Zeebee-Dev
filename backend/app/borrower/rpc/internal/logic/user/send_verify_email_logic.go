package userlogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
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
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	email := in.Data.Email
	// check send email limit
	key := stringx.Format(constants.RegisterVerifyEmail, constants.UserTypeBorrower, email)
	limitOpt := limiter.WithRedisLimit(
		global.Config.LimiterConf.SendVerityEmail.Rate,
		time.Second*time.Duration(global.Config.LimiterConf.SendVerityEmail.Period),
	)
	ok, err := global.SendRegisterEmailLimiter.Limit(l.ctx, key, limitOpt)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !ok {
		return nil, errorx.Errorf("Email %s Try Too Many Times", email)
	}

	// check user status
	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, email, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	if user.Status != constants.UserStatusUnverified {
		return nil, errorx.WithStack(errorx.New("User Is Verified"))
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

	// err = sendRegisterVerifyUrlEmail(l.ctx, email, user.Borrower.FirstName)
	// if err != nil {
	// 	return nil, err
	// }

	return &borrower.SendVerifyEmailResponse{}, nil
}
