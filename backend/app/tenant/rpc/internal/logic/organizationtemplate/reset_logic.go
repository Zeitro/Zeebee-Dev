package organizationtemplatelogic

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

type ResetLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Reset reverts a single organization template back to its base template.
func (l *ResetLogic) Reset(in *tenant.OrganizationTemplateResetRequest) (resp *tenant.OrganizationTemplateDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	result, err := l.baseLogic.OrganizationTemplateLogic.Reset(l.ctx, in.Data.Id, in.Data.OrganizationId)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.OrganizationTemplateEntityDto, proto.OrganizationTemplateEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationTemplateDetailResponse{Data: data}, nil
}
