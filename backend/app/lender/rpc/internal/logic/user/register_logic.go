package userlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"

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
		return nil, errorx.New("User Already Exists")
	}

	// check company
	isExisted, err := l.svcCtx.OrganizationRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"id": in.Data.OrganizationId,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !isExisted {
		return nil, errorx.New("Organization Is Not Exists")
	}

	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	userModel := &model.User{
		Email:        in.Data.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusActive,
		Type:         constants.UserTypeLender,
		FirstName:    in.Data.FirstName,
		LastName:     in.Data.LastName,
		FullName:     stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
	}
	// lenderEmployeeModel := &model.LenderEmployee{
	// 	FirstName:              in.Data.FirstName,
	// 	LastName:               in.Data.LastName,
	// 	FullName:               stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
	// 	DefaultLenderCompanyId: in.Data.LenderCompanyId,
	// }
	err = l.svcCtx.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var organization *model.Organization
		organization, innerErr = l.svcCtx.OrganizationRepo.FindOneByIdWithTx(l.ctx, tx, in.Data.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			return errorx.New("Organization Is Not Active")
		}
		userModel.Organizations = []*model.Organization{
			organization,
		}

		// create user
		userModel, innerErr = l.svcCtx.UserRepo.CreateWithTx(l.ctx, tx, userModel, true)
		if innerErr != nil {
			return innerErr
		}
		// create lender employee
		// lenderEmployeeModel.UserId = userModel.Id
		// lenderEmployeeModel.LenderCompanies = []*model.LenderCompany{
		// 	company,
		// }
		// _, innerErr = l.svcCtx.LenderEmployeeRepo.CreateEmployeeWithTx(l.ctx, tx, lenderEmployeeModel)
		// if innerErr != nil {
		// 	return innerErr
		// }

		return innerErr
	})
	if err != nil {
		return nil, err
	}
	// userModel.Organizations = lenderEmployeeModel
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
