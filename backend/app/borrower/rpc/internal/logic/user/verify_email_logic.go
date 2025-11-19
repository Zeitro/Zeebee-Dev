package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyEmailLogic) VerifyEmail(in *borrower.VerifyEmailRequest) (resp *borrower.VerifyEmailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var token string
	var user *model.User
	var jwtInfo *borrower.UserLoginResponseData

	email := in.Data.Email
	code := in.Data.Code
	if email == "" || code == "" {
		return nil, errorx.WithStack(errorx.New("Email And Code Are Required"))
	}
	// check code
	verifyCodeKey := stringx.Format(constants.BorrowerRegisterVerifyCode, email)
	token, err = global.GetRedis().Get(l.ctx, verifyCodeKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if token != code {
		return nil, errorx.New("Invalid Or Expired Code")
	}

	// get user
	user, err = l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, email)
	if err != nil {
		return nil, err
	}
	if user.Status != constants.UserStatusUnverified {
		return nil, errorx.New("Email Is Already Verified")
	}

	err = l.svcCtx.UserRepo.UpdateUserStatus(l.ctx, user.Id, constants.UserStatusActive)
	if err != nil {
		return nil, err
	}
	// generate jwt info
	jwtInfo, err = generateBorrowerToken(user, l.svcCtx)
	if err != nil {
		return nil, err
	}
	// remove code
	_ = global.GetRedis().Del(l.ctx, verifyCodeKey)

	return &borrower.VerifyEmailResponse{
		Data: jwtInfo,
	}, nil
}
