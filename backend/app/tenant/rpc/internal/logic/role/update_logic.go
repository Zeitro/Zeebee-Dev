package rolelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Update edits a role and its resource bindings.
func (l *UpdateLogic) Update(in *tenant.RoleUpdateRequest) (resp *tenant.RoleDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.RoleUpdateRequestData, dto.RoleUpdateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.Type = userType
	input.OrganizationId = in.Data.OrganizationId

	// before, _ := l.baseLogic.RoleLogic.Detail(l.ctx, in.Data.Id, in.Data.OrganizationId, userType)
	role, err := l.baseLogic.RoleLogic.Update(l.ctx, in.Data.Id, *input)
	if err != nil {
		return nil, err
	}
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeRole, role.Id, constants.OperationActionRoleUpdate, logic.AuditChange{Before: before, After: role})

	data, err := copierx.Trans[dto.RoleEntityDto, proto.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.RoleDetailResponse{Data: data}, nil
}
