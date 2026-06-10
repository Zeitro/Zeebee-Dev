package userlogic

import (
	"context"
	"fmt"

	"github.com/duke-git/lancet/v2/random"

	"zeebee-backend/app/lender/rpc/internal/email_tpl"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *RegisterLogic) Register(in *lender.UserRegisterRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	// check if user already exists
	var existedUser *model.User
	existedUser, _ = l.svcCtx.UserRepo.GetUserByEmail(l.ctx, in.Data.Email, constants.UserTypeLender)

	if existedUser != nil && existedUser.Id > 0 {
		ctx := context.WithoutCancel(l.ctx)
		if existedUser.Status == constants.UserStatusUnverified {
			// send verify email
			go func() {
				err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, existedUser.FirstName)
				if err != nil {
					err = errorx.Wrap(err, "Send Verify Email Error")
					l.WithContext(ctx).Error(err)
				}
			}()
		} else {
			// send login email
			go func() {
				err := sendLoginUrlEmail(ctx, in.Data.Email, existedUser.FirstName)
				if err != nil {
					err = errorx.Wrap(err, "Send Login Email Error")
					l.WithContext(ctx).Error(err)
				}
			}()
		}

		return &lender.EmptyResponse{}, nil
	}

	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	userInput := dto.UserRegisterRequestDto{
		Email:        in.Data.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusUnverified,
		Type:         constants.UserTypeLender,
		FirstName:    in.Data.FirstName,
		LastName:     in.Data.LastName,
		FullName:     stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
	}

	organizationInput := &dto.OrganizationRegisterRequestDto{
		Name:      in.Data.OrganizationName,
		Status:    constants.OrganizationStatusActive,
		OwnerType: constants.OrganizationOwnerTypeThirdParty,
		Type:      constants.UserTypeLender,
	}

	userModel, _, err := l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, userInput, *organizationInput, constants.UserTypeLender, false)
	if err != nil {
		return nil, err
	}

	newCxt := context.WithoutCancel(l.ctx)
	go func() {
		// send account validate email
		ctx := context.WithoutCancel(newCxt)
		err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, userModel.FirstName)
		if err != nil {
			l.WithContext(ctx).Error(err)
		}
		if err != nil {
			err = errorx.Wrap(err, "Send Register Verify Email Error")
			l.WithContext(ctx).Error(err)
		}
	}()

	return &lender.EmptyResponse{}, err
}

func sendLoginUrlEmail(ctx context.Context, email string, firstName string) (err error) {
	key := stringx.Format(constants.LoginEmail, constants.UserTypeLender, email)
	ok, err := global.LimiterMgr.Limit(ctx, limiter.StrategySendLoginEmail, key)

	if err != nil {
		return errorx.WithStack(err)
	}
	if !ok {
		return errorx.WithStack(errorx.Errorf("Email %s Try Too Many Times", email))
	}

	// send login email
	data := map[string]interface{}{
		"url":        global.Config.LoginUrl,
		"first_name": firstName,
	}
	content, err := email2.ParseTemplate(email_tpl.LoginUrlEmail, data)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := email2.Message{
		Subject: "Login to your account",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errorx.WithStack(err)
}

func sendRegisterVerifyUrlEmail(ctx context.Context, email string, firstName string) (err error) {
	var code string
	// generate verify code
	codeKey := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeLender, email)
	code, _ = global.GetRedis().Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		err = global.GetRedis().SetExpire(ctx, codeKey, code, global.Config.VerifyEmailExpireSeconds)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	// send verify email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s?code=%s&email=%s", global.Config.VerifyEmailUrl, code, email),
		"first_name": firstName,
	}
	content, err := email2.ParseTemplate(email_tpl.RegisterVerifyEmail, data)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := email2.Message{
		Subject: "Verify your email",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errorx.WithStack(err)
}
