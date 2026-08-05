package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// VerifyEmail confirms the email code, activates the user, and issues a TokenSet.
// tenant handles identity only; platform-specific side effects (e.g. borrower
// pre-approval letter) stay at the caller backend.
func (l *VerifyEmailLogic) VerifyEmail(in *tenant.VerifyEmailRequest) (resp *proto.UserLoginResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	user, err := l.baseLogic.UserLogic.VerifyUserEmail(l.ctx, in.Data.Email, in.Data.Code, userType)
	if err != nil {
		return nil, err
	}
	orgUser, err := l.svcCtx.OrganizationUserRepo.GetOrganizationUserByUserId(l.ctx, user.DefaultOrganizationId, user.Id)
	if err != nil {
		return nil, err
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, orgUser, userType)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}
	return &proto.UserLoginResponse{Data: data}, nil
}
