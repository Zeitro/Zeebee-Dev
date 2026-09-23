package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(in *admin.ProjectManagementTemplateDetailRequest) (resp *admin.ProjectManagementTemplateDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	template, err := base.ProjectManagementTemplateLogic.Detail(l.ctx, claim, in.Data.Id)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementTemplate, proto.ProjectManagementTemplateEntity](template, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	stages, err := copierx.Trans[[]*model.ProjectManagementStage, []*proto.ProjectManagementStageEntity](&template.Stages, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	data.Stages = *stages
	return &admin.ProjectManagementTemplateDetailResponse{Data: data}, nil
}
