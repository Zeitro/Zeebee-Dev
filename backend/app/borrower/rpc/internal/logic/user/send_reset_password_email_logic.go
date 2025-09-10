package userlogic

import (
	"context"
	"time"

	"github.com/duke-git/lancet/v2/random"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
	"zeebee-backend/common/stringx"
)

type SendResetPasswordEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendResetPasswordEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetPasswordEmailLogic {
	return &SendResetPasswordEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendResetPasswordEmailLogic) SendResetPasswordEmail(in *borrower.SendResetPasswordEmailRequest) (resp *borrower.SendResetPasswordEmailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	email := in.Data.Email
	if email == "" {
		return nil, errors.WithStack(errors.New("Email is required"))
	}

	// rate limit
	key := stringx.Format(constants.BorrowerResetPwdEmail, email)
	limitOpt := limiter.WithRedisLimit(
		global.Config.LimiterConf.SendResetPwdEmail.Rate,
		time.Second*time.Duration(global.Config.LimiterConf.SendResetPwdEmail.Period),
	)
	ok, err := global.SendResetPasswordEmailLimiter.Limit(l.ctx, key, limitOpt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, errors.Errorf("Email %s try too many times", email)
	}

	// get user
	user, err := l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, email)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if user == nil || user.Borrower == nil {
		return nil, errors.WithStack(errors.New("User not found"))
	}

	resetCode := random.RandNumeral(10)

	// send email
	err = l.sendResetPwdEmail(resetCode, email)

	// save reset code to redis
	redisKey := stringx.Format(constants.BorrowerResetPwdCode, email)
	expireSeconds := global.Config.ResetPwdExpireSeconds
	err = global.GetRedis().SetExpire(l.ctx, redisKey, resetCode, expireSeconds)
	if err != nil {
		return nil, errors.New("Failed to cache reset code")
	}

	return &borrower.SendResetPasswordEmailResponse{}, nil
}

func (l *SendResetPasswordEmailLogic) sendResetPwdEmail(code string, email string) (err error) {
	// todo send reset pwd email

	return err
}
