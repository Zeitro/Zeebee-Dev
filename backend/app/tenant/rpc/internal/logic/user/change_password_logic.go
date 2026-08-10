package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// ChangePassword verifies the old password and sets the new one. user_id and
// user type come from the caller's JWT claim; tenant is the identity authority.
func (l *ChangePasswordLogic) ChangePassword(in *tenant.ChangePasswordRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.New("request data is nil")
	}
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input := dto.UserChangePasswordRequestDto{
		UserId:      claim.UserId,
		UserType:    claim.UserType,
		OldPassword: in.Data.OldPassword,
		NewPassword: in.Data.NewPassword,
	}

	if err = l.baseLogic.UserLogic.ChangePassword(l.ctx, input); err != nil {
		return nil, err
	}
	l.baseLogic.RecordAudit(l.ctx, dto.OperationCreateDto{
		EntityType:   constants.OperationEntityTypeUser,
		EntityId:     claim.UserId,
		OperatorId:   claim.UserId,
		OperatorType: claim.UserType,
		Operation:    constants.OperationActionUserPasswordEdit,
	})
	return &proto.EmptyResponse{}, nil
}
