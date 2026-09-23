package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"
)

type DefaultsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDefaultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DefaultsLogic {
	return &DefaultsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DefaultsLogic) Defaults(in *admin.ProjectManagementTemplateDefaultsRequest) (resp *admin.ProjectManagementTemplateDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[admin.ProjectManagementTemplateDefaultsRequestData, dto.ProjectManagementTemplateUpdateDto](in.Data, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	template, err := base.ProjectManagementTemplateLogic.SetDefaults(l.ctx, claim, in.Data.Id, *input)
	if err != nil {
		return nil, err
	}
	return templateResponse(template)
}
