package userlogic

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *borrower.UserLoginRequest) (*borrower.UserLoginResponse, error) {
	// todo check input wrong pwd times
	user, err := l.svcCtx.UserRepo.FindOneByColumns(l.ctx,
		map[string]interface{}{
			"email": in.Email,
			"type":  constants.UserTypeBorrower,
		},
	)
	// check pwd
	if !utils.CheckPassword(in.Password, user.PasswordHash) {
		// todo add input wrong pwd times to redis
		return nil, errors.New("password is incorrect")
	}

	// check status
	if user.Status != constants.UserStatusActive {
		return nil, errors.New("user is not active")
	}
	// generate token
	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)), user.Email, "", l.svcCtx.Config.JwtAuth.AccessSecret, l.svcCtx.Config.JwtAuth.AccessExpire)
	if err != nil {
		return nil, err
	}

	return &borrower.UserLoginResponse{
		Token: token,
	}, nil
}
