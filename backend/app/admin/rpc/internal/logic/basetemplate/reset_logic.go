package basetemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type ResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResetLogic) Reset(in *proto.IdRequest) (resp *proto.EmptyResponse, err error) {

	id := in.Data
	_, err = l.svcCtx.TenantClient.OrganizationTemplate.ResetByBaseTemplateId(l.ctx, &tenant.OrganizationTemplateResetByBaseTemplateIdRequest{Data: &tenant.OrganizationTemplateResetByBaseTemplateIdRequestData{
		UserType:       string(constants.UserTypeAdmin),
		BaseTemplateId: id,
	}})
	return &proto.EmptyResponse{}, err
}
