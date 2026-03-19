package rolelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type AssignOrganizationUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAssignOrganizationUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignOrganizationUsersLogic {
	return &AssignOrganizationUsersLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AssignOrganizationUsersLogic) AssignOrganizationUsers(in *lender.RoleAssignOrganizationUsersRequest) (resp *lender.RoleDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[lender.RoleAssignOrganizationUsersRequestData, dto.RoleAssignOrganizationUsersRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.Type = constants.UserTypeLender
	input.OrganizationId = claim.OrganizationId

	role, err := l.baseLogic.RoleLogic.AssignOrganizationUsers(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.RoleEntityDto, lender.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &lender.RoleDetailResponse{
		Data: data,
	}, nil
}
