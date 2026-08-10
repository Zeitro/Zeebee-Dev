package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Delete removes an organization member (org + platform scoped).
func (l *DeleteLogic) Delete(in *tenant.OrganizationUserDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	if err = l.baseLogic.OrganizationUserLogic.Delete(l.ctx, in.Data.Id, in.Data.OrganizationId, userType); err != nil {
		return
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationUser, in.Data.Id, constants.OperationActionMemberDelete)
	return &proto.EmptyResponse{}, nil
}
