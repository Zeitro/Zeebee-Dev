package userlogic

import (
	"context"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *borrower.UserRegisterRequest) (resp *borrower.UserRegisterResponse, err error) {

	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// check if user already exists
	var existedUser *model.User
	existedUser, _ = l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, in.Data.Email)

	if existedUser != nil && existedUser.Id > 0 {
		ctx := context.WithoutCancel(l.ctx)
		if existedUser.Status == constants.UserStatusUnverified {
			// send verify email
			go func() {
				err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, existedUser.Borrower.FirstName)
				if err != nil {
					err = errors.Wrap(err, "send verify email error")
					l.WithContext(ctx).Error(err)
				}
			}()
		} else {
			// send login email
			go func() {
				err := sendLoginUrlEmail(ctx, in.Data.Email, existedUser.Borrower.FirstName)
				if err != nil {
					err = errors.Wrap(err, "send login email error")
					l.WithContext(ctx).Error(err)
				}
			}()
		}

		return &borrower.UserRegisterResponse{
			Data: &borrower.UserRegisterResponseData{},
		}, nil
	}
	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userModel := &model.User{
		Email:        in.Data.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusUnverified,
		Type:         constants.UserTypeBorrower,
	}
	borrowerModel := &model.Borrower{
		FirstName: in.Data.FirstName,
		LastName:  in.Data.LastName,
		FullName:  stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
	}
	err = l.svcCtx.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// create user
		userModel, err = l.svcCtx.UserRepo.CreateWithTx(l.ctx, tx, userModel, false)
		if err != nil {
			return err
		}
		// create borrower
		borrowerModel.UserId = userModel.Id
		_, err = l.svcCtx.BorrowerRepo.CreateWithTx(l.ctx, tx, borrowerModel, false)
		return err
	})
	if err != nil {
		return nil, err
	}

	go func() {
		// send account validate email
		ctx := context.WithoutCancel(l.ctx)
		err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, borrowerModel.FirstName)
		if err != nil {
			l.WithContext(ctx).Error(err)
		}
		if err != nil {
			err = errors.Wrap(err, "send register verify email error")
			l.WithContext(ctx).Error(err)
		}
	}()

	return &borrower.UserRegisterResponse{
		Data: &borrower.UserRegisterResponseData{},
	}, err
}

// sendLoginUrlEmail send login url email
func sendLoginUrlEmail(ctx context.Context, email string, firstName string) (err error) {
	key := stringx.Format(constants.BorrowerLoginEmail, email)
	limitOpt := limiter.WithRedisLimit(
		global.Config.LimiterConf.SendLoginEmail.Rate,
		time.Second*time.Duration(global.Config.LimiterConf.SendLoginEmail.Period),
	)
	ok, err := global.SendLoginEmailLimiter.Limit(ctx, key, limitOpt)
	if err != nil {
		return errors.WithStack(err)
	}
	if !ok {
		return errors.WithStack(errors.Errorf("Email %s try too many times", email))
	}

	// send login email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s", global.Config.LoginUrl),
		"first_name": firstName,
	}
	content, err := email_tpl.GetParsedContent(email_tpl.LoginUrlEmail, data)
	if err != nil {
		return errors.WithStack(err)
	}
	msg := email2.EmailMessage{
		Subject: "Login to your account",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errors.WithStack(err)
}

// sendRegisterVerifyUrlEmail send register verify url email
func sendRegisterVerifyUrlEmail(ctx context.Context, email string, firstName string) (err error) {
	var code string
	// generate verify code
	codeKey := stringx.Format(constants.BorrowerRegisterVerifyCode, email)
	code, _ = global.GetRedis().Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		err = global.GetRedis().SetExpire(ctx, codeKey, code, global.Config.VerifyEmailExpireSeconds)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	// send verify email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s?code=%s&email=%s", global.Config.VerifyEmailUrl, code, email),
		"first_name": firstName,
	}
	content, err := email_tpl.GetParsedContent(email_tpl.RegisterVerifyEmail, data)
	if err != nil {
		return errors.WithStack(err)
	}
	msg := email2.EmailMessage{
		Subject: "Verify your email",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errors.WithStack(err)
}
