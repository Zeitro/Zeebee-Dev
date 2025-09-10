package userlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/utils"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *borrower.ResetPasswordRequest) (resp *borrower.ResetPasswordResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var jwtInfo *borrower.UserLoginResponseData

	email := in.Data.Email
	code := in.Data.Code
	newPwd := in.Data.Password
	if email == "" || code == "" || newPwd == "" {
		return nil, errors.WithStack(errors.New("Email, code and new password are required"))
	}
	redisKey := stringx.Format(constants.BorrowerResetPwdCode, email)
	cacheToken, err := global.GetRedis().Get(l.ctx, redisKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if cacheToken == "" || cacheToken != code {
		return nil, errors.WithStack(errors.New("Invalid or expired code"))
	}

	user, err := l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, email)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if user.Borrower == nil {
		return nil, errors.WithStack(errors.New("User not found"))
	}

	// hash pwd
	pwdHash, err := utils.HashPassword(newPwd)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// update pwd
	err = l.svcCtx.UserRepo.UpdateUserPassword(l.ctx, user.Id, pwdHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	jwtInfo, err = generateBorrowerToken(user, l.svcCtx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// delete redis key
	_ = global.GetRedis().Del(l.ctx, redisKey)
	return &borrower.ResetPasswordResponse{
		Data: jwtInfo,
	}, nil
}
