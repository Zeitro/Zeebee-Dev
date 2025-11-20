package userlogic

import (
	"context"

	"zeebee-backend/domains/zeebee/model"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchLenderCompanyLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewSwitchLenderCompanyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchLenderCompanyLogic {
	return &SwitchLenderCompanyLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *SwitchLenderCompanyLogic) SwitchLenderCompany(in *lender.SwitchLenderCompanyRequest) (resp *lender.RefreshTokenResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// get  employee
	employee, err := l.svcCtx.LenderEmployeeRepo.GetLenderEmployee(l.ctx, in.Data.EmployeeId)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// check company
	isValidCompany := false
	slice.ForEach(employee.LenderCompanies, func(index int, item *model.LenderCompany) {
		if item.Id == in.Data.LenderCompanyId && item.Status == constants.LenderCompanyStatusActive {
			isValidCompany = true
		}
	})

	if !isValidCompany {
		return nil, errorx.New("Invalid Lender Company For This Employee")
	}

	// update employee
	employee.DefaultLenderCompanyId = in.Data.LenderCompanyId
	err = l.svcCtx.LenderEmployeeRepo.UpdateColumnsById(l.ctx, employee.Id, map[string]interface{}{
		"default_lender_company_id": employee.DefaultLenderCompanyId,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// get user info
	user, err := l.svcCtx.UserRepo.GetLenderUser(l.ctx, employee.UserId)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	if user.Status == constants.UserStatusInactive {
		return nil, errorx.New("User Is Inactive")
	}

	// generate token
	token, err := generateLenderToken(user, l.svcCtx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &lender.RefreshTokenResponse{
		Data: &lender.RefreshTokenResponseData{
			AccessToken:        token.AccessToken,
			AccessTokenExpire:  token.AccessTokenExpire,
			RefreshToken:       token.RefreshToken,
			RefreshTokenExpire: token.RefreshTokenExpire,
		},
	}, nil
}
