package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/errorx"
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
	page, err := logic.GetBaseLogic(l.svcCtx).ProjectManagementTemplateLogic.Page(l.ctx, in.Data.Page, in.Data.Limit)
	if err != nil {
		return nil, err
	}
	list := make([]*proto.ProjectManagementTemplateEntity, 0, len(page.List))
	for _, item := range page.List {
		entity, convertErr := templateEntity(item)
		if convertErr != nil {
			return nil, convertErr
		}
		list = append(list, entity)
	}
	return &admin.ProjectManagementTemplatePageResponse{Total: page.Total, Page: page.Page, Limit: page.Limit, TotalPage: page.TotalPage, List: list}, nil
}
