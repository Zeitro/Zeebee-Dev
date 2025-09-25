package userlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *lender.UserInfoRequest) (resp *lender.UserinfoResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	user, err := l.svcCtx.UserRepo.GetLenderUser(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}
	userLender := user.LenderEmployee
	if userLender == nil {
		return nil, errors.WithStack(errors.New("lender not found"))
	}

	// Get lender company info
	var lenderCompanyId int64
	var lenderCompanyName string
	if userLender.LenderCompany != nil {
		lenderCompanyId = userLender.LenderCompanyId
		lenderCompanyName = userLender.LenderCompany.Name
	}

	userInfo := &lender.UserinfoResponse{
		Data: &lender.UserinfoResponseData{
			Id:                user.Id,
			Email:             user.Email,
			FirstName:         userLender.FirstName,
			LastName:          userLender.LastName,
			FullName:          userLender.FullName,
			LenderCompanyId:   lenderCompanyId,
			LenderCompanyName: lenderCompanyName,
			Status:            string(user.Status),
		},
	}
	return userInfo, nil
}
