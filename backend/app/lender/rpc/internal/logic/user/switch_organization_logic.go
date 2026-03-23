package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

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
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(in *lender.SwitchOrganizationRequest) (resp *lender.UserJwtResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	user, err := l.baseLogic.UserLogic.SwitchOrganization(l.ctx, claim.UserId, in.Data.OrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	// token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser, constants.UserTypeLender)
	token, err := l.baseLogic.UserLogic.GenerateJwtTokenById(l.ctx, user.Id, user.DefaultOrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.UserJwtTokenDto, lender.UserJwtResponseData](token)
	if err != nil {
		return nil, err
	}

	return &lender.UserJwtResponse{
		Data: data,
	}, nil
}
