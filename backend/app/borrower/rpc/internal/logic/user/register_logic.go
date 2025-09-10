package userlogic

import (
	"context"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
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
			l.WithContext(l.ctx).Error(err)
		}
	}()
	var isExisted bool
	isExisted, err = l.svcCtx.UserRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"email": in.Data.Email,
	})
	if err != nil {
		return nil, err
	}
	if isExisted {
		// send login email
		err = sendLoginUrlEmail(l.ctx, in.Data.Email)
		return nil, err
	}

	pwd, err := utils.HashPassword(in.Data.Password)
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
		// Ssn:         in.Data.Ssn,
		// PhoneNumber: in.Data.PhoneNumber,
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
		return nil
	})
	if err != nil {
		return nil, err
	}

	// send account validate email
	err = sendRegisterVerifyUrlEmail(l.ctx, in.Data.Email)

	return &borrower.UserRegisterResponse{
		Data: &borrower.UserRegisterResponseData{
			// AccessToken:        accessToken,
			// RefreshToken:       refreshToken,
			// AccessTokenExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
			// RefreshTokenExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
		},
	}, errors.WithStack(err)
}

// sendLoginUrlEmail send login url email
func sendLoginUrlEmail(ctx context.Context, email string) (err error) {
	key := stringx.Format(constants.BorrowerLoginEmail, email)
	limitOpt := limiter.WithRedisLimit(
		global.Config.LimiterConf.SendLoginEmail.Rate,
		time.Second*time.Duration(global.Config.LimiterConf.SendLoginEmail.Period),
	)
	ok, err := global.SendLoginEmailLimiter.Limit(ctx, key, limitOpt)
	if err != nil {
		return err
	}
	if !ok {
		return errors.WithStack(errors.Errorf("Email %s try too many times", email))
	}

	// todo send login email

	return err
}

// sendRegisterVerifyUrlEmail send register verify url email
func sendRegisterVerifyUrlEmail(ctx context.Context, email string) (err error) {
	// generate verify code
	codeKey := stringx.Format(constants.BorrowerRegisterVerifyCode, email)
	code := random.RandNumeral(10)
	err = global.GetRedis().SetExpire(ctx, codeKey, code, global.Config.VerifyEmailExpireSeconds)

	// todo send verify email

	return err
}
