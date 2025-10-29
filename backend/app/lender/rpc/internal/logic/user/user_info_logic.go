package userlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
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
	companies := []*lender.UserinfoResponseData_Company{}

	user, err := l.svcCtx.UserRepo.GetLenderUser(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}
	userLender := user.LenderEmployee
	if userLender == nil {
		return nil, errors.WithStack(errors.New("lender not found"))
	}
	if len(userLender.LenderCompanies) > 0 {
		err = copier.CopyWithOption(&companies, userLender.LenderCompanies, copier.Option{
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to copy lender companies")
		}
	}
	var currentCompany *lender.UserinfoResponseData_Company
	slice.ForEach(companies, func(i int, company *lender.UserinfoResponseData_Company) {
		if company.Id == userLender.DefaultLenderCompanyId {
			currentCompany = companies[i]
			company.IsDefault = true
		}
	})

	// Get lender company info
	var lenderCompanyId int64
	var lenderCompanyName string
	if currentCompany != nil {
		lenderCompanyId = userLender.DefaultLenderCompanyId
		lenderCompanyName = currentCompany.Name
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
			LenderCompanies:   companies,
		},
	}
	return userInfo, nil
}
