package userlogic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *borrower.UserInfoRequest) (resp *borrower.UserinfoResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).Error(err)
		}
	}()

	user, err := l.svcCtx.UserRepo.GetBorrowerUser(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}
	userBorrower := user.Borrower
	if userBorrower == nil {
		return nil, errors.WithStack(errors.New("borrower not found"))
	}
	userInfo := &borrower.UserinfoResponse{
		Data: &borrower.UserinfoResponseData{
			Id:          user.Id,
			Email:       user.Email,
			FirstName:   userBorrower.FirstName,
			LastName:    userBorrower.LastName,
			PhoneNumber: userBorrower.PhoneNumber,
			// Address:     userBorrower.Address,
			Ssn: userBorrower.Ssn,
		},
	}
	return userInfo, nil
}
