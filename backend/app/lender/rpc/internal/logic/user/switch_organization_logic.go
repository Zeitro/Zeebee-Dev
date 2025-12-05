package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	logic2 "zeebee-backend/domains/zeebee/logic"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SwitchOrganizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(in *lender.SwitchOrganizationRequest) (resp *lender.RefreshTokenResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	user, organizationUser, err := l.baseLogic.UserLogic.SwitchOrganization(l.ctx, in.Data.UserId, in.Data.OrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	jwtConfig := logic2.JwtConfig{
		AccessSecret:  l.svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: l.svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
	}
	token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser, jwtConfig, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data := &lender.RefreshTokenResponseData{}
	err = copier.CopyWithOption(data, token, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &lender.RefreshTokenResponse{Data: data}, nil
}
