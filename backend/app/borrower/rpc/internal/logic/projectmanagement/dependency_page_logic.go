package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type DependencyPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDependencyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DependencyPageLogic {
	return &DependencyPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DependencyPageLogic) DependencyPage(in *borrower.ProjectManagementDependencyPageRequest) (resp *borrower.ProjectManagementDependencyPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	page, err := base.ProjectManagementLogic.PageTaskDependencies(l.ctx, actor, dto.ProjectManagementTaskDependencyPageRequestDto{TaskId: in.Data.TaskId, Page: in.Data.Page, Limit: in.Data.Limit, Predecessor: in.Data.Predecessor})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data := &borrower.ProjectManagementDependencyPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*borrower.ProjectManagementTaskDependencyEntity, 0, len(page.List))}
	for _, dependency := range page.List {
		entity, err := dependencyEntity(dependency)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementDependencyPageResponse{Data: data}, nil
}
