package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/constants"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type StatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusLogic {
	return &StatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatusLogic) Status(in *admin.ProjectManagementTemplateStatusRequest) (resp *admin.ProjectManagementTemplateDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	template, err := base.ProjectManagementTemplateLogic.SetStatus(l.ctx, claim, in.Data.Id, constants.CommonStatus(in.Data.Status))
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementTemplate, proto.ProjectManagementTemplateEntity](template, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementTemplateDetailResponse{Data: data}, nil
}
