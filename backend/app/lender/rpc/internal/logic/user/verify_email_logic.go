package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *VerifyEmailLogic) VerifyEmail(in *lender.VerifyEmailRequest) (resp *lender.UserJwtResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var user *model.User
	email := in.Data.Email
	code := in.Data.Code
	user, err = l.baseLogic.UserLogic.VerifyUserEmail(l.ctx, email, code, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	organizationUser, err := l.svcCtx.OrganizationUserRepo.GetOrganizationUserByUserId(l.ctx, user.DefaultOrganizationId, user.Id)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserJwtTokenDto, lender.UserJwtResponseData](tokenInfo)
	if err != nil {
		return nil, err
	}
	return &lender.UserJwtResponse{
		Data: data,
	}, nil
}
