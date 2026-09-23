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

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(in *admin.ProjectManagementTemplatePageRequest) (resp *admin.ProjectManagementTemplatePageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	pageReq, err := copierx.Trans[admin.ProjectManagementTemplatePageRequestData, dto.ProjectManagementTemplatePageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	page, err := logic.GetBaseLogic(l.svcCtx).ProjectManagementTemplateLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*model.ProjectManagementTemplate, []*proto.ProjectManagementTemplateEntity](&page.List, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementTemplatePageResponse{Total: page.Total, Page: page.Page, Limit: page.Limit, TotalPage: page.TotalPage, List: *list}, nil
}
