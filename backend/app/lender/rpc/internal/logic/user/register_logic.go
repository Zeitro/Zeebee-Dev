package userlogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
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

func (l *RegisterLogic) Register(in *lender.UserRegisterRequest) (resp *lender.UserRegisterResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// check if user already exists
	var existedUser *model.User
	existedUser, _ = l.svcCtx.UserRepo.GetLenderUserByEmail(l.ctx, in.Data.Email)

	if existedUser != nil && existedUser.Id > 0 {
		return nil, errors.WithStack(errors.New("User already exists"))
	}

	// check company
	isExisted, err := l.svcCtx.LenderCompanyRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"id": in.Data.LenderCompanyId,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !isExisted {
		return nil, errors.WithStack(errors.New("Lender company is not exists"))
	}

	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userModel := &model.User{
		Email:        in.Data.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusActive,
		Type:         constants.UserTypeLender,
	}
	lenderEmployeeModel := &model.LenderEmployee{
		FirstName:       in.Data.FirstName,
		LastName:        in.Data.LastName,
		FullName:        stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
		LenderCompanyId: in.Data.LenderCompanyId,
	}
	err = l.svcCtx.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// create user
		userModel, err = l.svcCtx.UserRepo.CreateWithTx(l.ctx, tx, userModel, false)
		if err != nil {
			return err
		}
		// create lender employee
		lenderEmployeeModel.UserId = userModel.Id
		_, err = l.svcCtx.LenderEmployeeRepo.CreateWithTx(l.ctx, tx, lenderEmployeeModel, false)
		return err
	})
	if err != nil {
		return nil, err
	}
	userModel.LenderEmployee = lenderEmployeeModel
	token, err := generateLenderToken(userModel, l.svcCtx)
	if err != nil {
		return nil, err
	}

	return &lender.UserRegisterResponse{
		Data: &lender.UserRegisterResponseData{
			AccessToken:        token.AccessToken,
			RefreshToken:       token.RefreshToken,
			AccessTokenExpire:  token.AccessTokenExpire,
			RefreshTokenExpire: token.RefreshTokenExpire,
		},
	}, err
}
