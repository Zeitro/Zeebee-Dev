package rolelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *AssignOrganizationUsersLogic) AssignOrganizationUsers(in *borrower.RoleAssignOrganizationUsersRequest) (resp *borrower.RoleDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.RoleAssignOrganizationUsersRequestData, dto.RoleAssignOrganizationUsersRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.Type = constants.UserTypeBorrower
	input.OrganizationId = claim.OrganizationId

	role, err := l.baseLogic.RoleLogic.AssignOrganizationUsers(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.RoleEntityDto, proto.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.RoleDetailResponse{
		Data: data,
	}, nil
}
