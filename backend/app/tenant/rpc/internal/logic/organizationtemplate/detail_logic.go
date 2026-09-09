package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Detail returns an organization template by id.
func (l *DetailLogic) Detail(in *tenant.OrganizationTemplateDetailRequest) (resp *tenant.OrganizationTemplateDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	result, err := l.baseLogic.OrganizationTemplateLogic.Detail(l.ctx, in.Data.Id)
	if err != nil {
		return
	}
	if userType != constants.UserTypeAdmin && result.OrganizationId != in.Data.OrganizationId {
		return nil, errorx.New("template does not belong to organization")
	}

	data, err := copierx.Trans[dto.OrganizationTemplateEntityDto, proto.OrganizationTemplateEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationTemplateDetailResponse{Data: data}, nil
}
