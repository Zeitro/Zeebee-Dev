package userlogic

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/limiter"
	"zeebee-backend/common/stringx"
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
			l.WithContext(l.ctx).Error(err)
		}
	}()

	email := in.Data.Email
	// check send email limit
	key := stringx.Format(constants.BorrowerRegisterVerifyEmail, email)
	limitOpt := limiter.WithRedisLimit(
		global.Config.LimiterConf.SendVerityEmail.Rate,
		time.Second*time.Duration(global.Config.LimiterConf.SendVerityEmail.Period),
	)
	ok, err := global.SendLoginEmailLimiter.Limit(l.ctx, key, limitOpt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, errors.WithStack(errors.Errorf("Email %s try too many times", email))
	}

	// check user status
	user, err := l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, email)
	if err != nil {
		return nil, err
	}

	if user.Status != constants.UserStatusUnverified {
		return nil, errors.WithStack(errors.New("User is verified"))
	}

	if user.Borrower == nil {
		return nil, errors.WithStack(errors.New("User is not found"))
	}

	// send email
	go func() {
		// send account validate email
		ctx := context.WithoutCancel(l.ctx)
		err := sendRegisterVerifyUrlEmail(ctx, email, user.Borrower.FirstName)
		if err != nil {
			l.WithContext(ctx).Error(err)
		}
		if err != nil {
			err = errors.Wrap(err, "send register verify email error")
			l.WithContext(ctx).Error(err)
		}
	}()

	// err = sendRegisterVerifyUrlEmail(l.ctx, email, user.Borrower.FirstName)
	// if err != nil {
	// 	return nil, err
	// }

	return &borrower.SendVerifyEmailResponse{}, nil
}
