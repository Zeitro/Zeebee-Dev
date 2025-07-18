package userlogic

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
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

func (l *RegisterLogic) Register(in *borrower.UserRegisterRequest) (*borrower.UserRegisterResponse, error) {

	isExisted, err := l.svcCtx.UserRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"email": in.Email,
	})
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errors.New("email already existed")
	}

	pwd, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	userModel := &model.User{
		Email:        in.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusActive,
		Type:         constants.UserTypeBorrower,
	}
	err = l.svcCtx.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// create user
		userModel, err = l.svcCtx.UserRepo.CreateWithTx(l.ctx, tx, userModel)
		if err != nil {
			return err
		}
		// create borrower
		borrowerModel := &model.Borrower{
			UserID:      userModel.ID,
			FirstName:   in.FirstName,
			LastName:    in.LastName,
			Ssn:         in.Ssn,
			PhoneNumber: in.PhoneNumber,
			Address:     in.Address,
		}
		_, err = l.svcCtx.BorrowerRepo.CreateWithTx(l.ctx, tx, borrowerModel)
		return nil
	})

	// generate token
	token, err := utils.GenerateToken(
		strconv.Itoa(int(userModel.ID)),
		userModel.Email,
		"",
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, err
	}

	return &borrower.UserRegisterResponse{
		Token: token,
	}, nil
}
