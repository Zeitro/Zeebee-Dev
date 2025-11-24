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

type SwitchOrganizationLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(in *lender.SwitchOrganizationRequest) (resp *lender.RefreshTokenResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	user, err := l.svcCtx.UserRepo.GetLenderUser(l.ctx, in.Data.UserId)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// check company
	isValidCompany := false
	slice.ForEach(user.Organizations, func(index int, item *model.Organization) {
		if item.Id == in.Data.OrganizationId && item.Status == constants.OrganizationStatusActive {
			isValidCompany = true
		}
	})

	if !isValidCompany {
		return nil, errorx.New("Invalid Organization For This Employee")
	}

	// update employee
	user.DefaultOrganizationId = in.Data.OrganizationId
	err = l.svcCtx.UserRepo.UpdateColumnsById(l.ctx, user.Id, map[string]interface{}{
		"default_organization_id": user.DefaultOrganizationId,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// get user info
	user, err = l.svcCtx.UserRepo.GetLenderUser(l.ctx, user.Id)
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
