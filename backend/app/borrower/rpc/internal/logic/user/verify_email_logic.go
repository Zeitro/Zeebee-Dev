package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	logic2 "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
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
	user, err = l.svcCtx.UserRepo.GetUserByEmail(l.ctx, email, constants.UserTypeBorrower)
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

	tokenInfo, err := l.baseLogic.UserLogic.GenerateJwtToken(user, logic2.JwtConfig{
		AccessSecret:  l.svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: l.svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
	}, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &borrower.UserLoginResponseData{}
	err = copier.CopyWithOption(data, tokenInfo, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.VerifyEmailResponse{
		Data: data,
	}, nil
}
