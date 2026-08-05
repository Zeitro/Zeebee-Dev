package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type ResetByBaseTemplateIdLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewResetByBaseTemplateIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetByBaseTemplateIdLogic {
	return &ResetByBaseTemplateIdLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// ResetByBaseTemplateId cascades a reset of all organization templates derived
// from the given base template (admin operation).
func (l *ResetByBaseTemplateIdLogic) ResetByBaseTemplateId(in *tenant.OrganizationTemplateResetByBaseTemplateIdRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	if err = l.baseLogic.OrganizationTemplateLogic.ResetByBaseTemplateId(l.ctx, in.Data.BaseTemplateId); err != nil {
		return
	}
	return &proto.EmptyResponse{}, nil
}
