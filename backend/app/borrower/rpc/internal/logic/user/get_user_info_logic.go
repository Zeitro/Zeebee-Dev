package userlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/gormx"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *GetUserInfoLogic) GetUserInfo(in *borrower.UserInfoRequest) (*borrower.UserinfoResponse, error) {

	opt := gormx.NewQueryOption(
		gormx.WithFilterOpt([]gormx.WhereRaw{{Raw: "id = ?", Args: []interface{}{in.Id}}}),
		gormx.WithPreloadsOpt([]string{"Borrower"}),
	)

	user, err := l.svcCtx.UserRepo.FindOne(l.ctx, opt)

	// opt := repo.QueryOpt{
	// 	Preloads: []string{"Borrower"},
	// }
	// user, err := repo.GetOneById[*model.User](l.ctx, in.Id, opt)
	if err != nil {
		return nil, err
	}
	userBorrower := user.Borrower
	if userBorrower == nil {
		return nil, errors.New("borrower not found")
	}
	userInfo := &borrower.UserinfoResponse{
		Id:          user.ID,
		Email:       user.Email,
		FirstName:   userBorrower.FirstName,
		LastName:    userBorrower.LastName,
		PhoneNumber: userBorrower.PhoneNumber,
		Address:     userBorrower.Address,
		Ssn:         userBorrower.Ssn,
	}
	return userInfo, nil
}
