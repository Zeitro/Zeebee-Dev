package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type ResetPasswordByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResetPasswordByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordByCodeLogic {
	return &ResetPasswordByCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// ResetPasswordByCode verifies the one-time code, sets the new password, and
// returns a fresh TokenSet signed with the platform key.
func (l *ResetPasswordByCodeLogic) ResetPasswordByCode(in *tenant.ResetPasswordByCodeRequest) (resp *proto.UserLoginResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.ResetPasswordByCodeRequestData, dto.UserResetPasswordByEmailRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = userType

	token, err := l.baseLogic.UserLogic.ResetPasswordWithEmail(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	l.baseLogic.RecordAudit(l.ctx, dto.OperationCreateDto{
		EntityType:   constants.OperationEntityTypeUser,
		OperatorName: in.Data.Email,
		OperatorType: userType,
		Operation:    constants.OperationActionUserPasswordReset,
	})

	data, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}
	return &proto.UserLoginResponse{Data: data}, nil
}
