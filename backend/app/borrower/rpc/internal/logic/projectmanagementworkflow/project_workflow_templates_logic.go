package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type ProjectWorkflowTemplatesLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectWorkflowTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowTemplatesLogic {
	return &ProjectWorkflowTemplatesLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectWorkflowTemplatesLogic) ProjectWorkflowTemplates(in *proto.EmptyRequest) (resp *borrower.ProjectManagementWorkflowTemplateResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	templates, err := l.baseLogic.ProjectManagementTemplateLogic.Discover(l.ctx, claim)
	if err != nil {
		return nil, err
	}
	list := make([]*proto.ProjectManagementTemplateEntity, 0, len(templates))
	for _, template := range templates {
		entity, convertErr := copierx.Trans[model.ProjectManagementTemplate, proto.ProjectManagementTemplateEntity](template, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
		if convertErr != nil {
			return nil, convertErr
		}
		list = append(list, entity)
	}
	return &borrower.ProjectManagementWorkflowTemplateResponse{List: list}, nil
}
