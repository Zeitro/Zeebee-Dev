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
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type CreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.ProjectManagementTemplateCreateRequest) (resp *admin.ProjectManagementTemplateDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[admin.ProjectManagementTemplateCreateRequestData, dto.ProjectManagementTemplateCreateDto](in.Data, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	template, err := base.ProjectManagementTemplateLogic.Create(l.ctx, claim, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementTemplate, proto.ProjectManagementTemplateEntity](template, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementTemplateDetailResponse{Data: data}, nil
}
