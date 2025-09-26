package userlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *ResetPasswordLogic) ResetPassword(in *lender.ResetPasswordRequest) (resp *lender.ResetPasswordResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var jwtInfo *lender.UserLoginResponseData

	email := in.Data.Email
	code := in.Data.Code
	newPwd := in.Data.Password
	if email == "" || code == "" || newPwd == "" {
		return nil, errors.WithStack(errors.New("Email, code and new password are required"))
	}
	redisKey := stringx.Format(constants.LenderResetPwdCode, email)
	cacheToken, err := global.GetRedis().Get(l.ctx, redisKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if cacheToken == "" || cacheToken != code {
		return nil, errors.WithStack(errors.New("Invalid or expired code"))
	}

	user, err := l.svcCtx.UserRepo.GetLenderUserByEmail(l.ctx, email)
	if err != nil {
		return nil, err
	}
	if user.LenderEmployee == nil {
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
		return nil, err
	}
	jwtInfo, err = generateLenderToken(user, l.svcCtx)
	if err != nil {
		return nil, err
	}
	// delete redis key
	_ = global.GetRedis().Del(l.ctx, redisKey)
	return &lender.ResetPasswordResponse{
		Data: jwtInfo,
	}, nil
}
