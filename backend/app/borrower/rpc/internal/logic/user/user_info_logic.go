package userlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
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

func (l *UserInfoLogic) UserInfo(in *borrower.UserInfoRequest) (resp *borrower.UserinfoResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	organizations := []*borrower.UserinfoResponseData_Organization{}

	user, err := l.svcCtx.UserRepo.GetUser(l.ctx, in.Data.Id, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User Not Found"))
	}
	if len(user.Organizations) > 0 {
		err = copier.CopyWithOption(&organizations, user.Organizations, copier.Option{
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if err != nil {
			return nil, errorx.Wrap(err, "Failed To Copy Organization")
		}
	}
	var currentOrganization *borrower.UserinfoResponseData_Organization
	slice.ForEach(organizations, func(i int, organization *borrower.UserinfoResponseData_Organization) {
		if organization.Id == user.DefaultOrganizationId {
			currentOrganization = organizations[i]
			organization.IsDefault = true
		}
	})

	// Get lender company info
	var lenderCompanyId int64
	var lenderCompanyName string
	if currentOrganization != nil {
		lenderCompanyId = user.DefaultOrganizationId
		lenderCompanyName = currentOrganization.Name
	}

	userInfo := &borrower.UserinfoResponse{
		Data: &borrower.UserinfoResponseData{
			Id:               user.Id,
			Email:            user.Email,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			FullName:         user.FullName,
			OrganizationId:   lenderCompanyId,
			OrganizationName: lenderCompanyName,
			Status:           string(user.Status),
			Organizations:    organizations,
		},
	}
	return userInfo, nil
}
