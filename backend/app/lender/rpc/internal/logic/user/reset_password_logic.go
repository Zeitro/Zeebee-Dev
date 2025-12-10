package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *ResetPasswordLogic) ResetPassword(in *lender.ResetPasswordRequest) (resp *lender.ResetPasswordResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var jwtInfo *lender.UserLoginResponseData = &lender.UserLoginResponseData{}
	var organizationUser *model.OrganizationUser

	email := in.Data.Email
	code := in.Data.Code
	newPwd := in.Data.Password
	if email == "" || code == "" || newPwd == "" {
		return nil, errorx.WithStack(errorx.New("Email, Code And New Password Are Required"))
	}
	redisKey := stringx.Format(constants.LenderResetPwdCode, email)
	cacheToken, err := global.GetRedis().Get(l.ctx, redisKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if cacheToken == "" || cacheToken != code {
		return nil, errorx.WithStack(errorx.New("Invalid Or Expired Code"))
	}

	user, err := l.svcCtx.UserRepo.GetLenderUserByEmail(l.ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User Not Found"))
	}

	// hash pwd
	pwdHash, err := utils.HashPassword(newPwd)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// update pwd
	err = l.svcCtx.UserRepo.UpdateUserPassword(l.ctx, user.Id, pwdHash)
	if err != nil {
		return nil, err
	}

	user, organizationUser, err = l.baseLogic.UserLogic.CheckAndUpdateActiveOrganizationByUser(l.ctx, user)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	err = copier.CopyWithOption(jwtInfo, token, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	// delete redis key
	_ = global.GetRedis().Del(l.ctx, redisKey)
	return &lender.ResetPasswordResponse{
		Data: jwtInfo,
	}, nil
}
