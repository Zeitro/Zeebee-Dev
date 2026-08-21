package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
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

func (l *PageLogic) Page(in *borrower.ProjectManagementTaskPageRequest) (resp *borrower.ProjectManagementTaskPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := claimValue
	filter := dto.ProjectManagementTaskFilterDto{ProjectId: in.Data.ProjectId}
	if in.Data.Filter != nil {
		filter.StageId, filter.TaskTypeId = in.Data.Filter.StageId, in.Data.Filter.TaskTypeId
		filter.AssigneeId, filter.Keyword, filter.Overdue = in.Data.Filter.AssigneeId, in.Data.Filter.Keyword, in.Data.Filter.Overdue
	}
	input := dto.ProjectManagementTaskPageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit, Filter: filter}
	page, err := base.ProjectManagementTaskLogic.Page(l.ctx, claim, input)
	if err != nil {
		return nil, err
	}
	data := &borrower.ProjectManagementTaskPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.ProjectManagementTaskEntity, 0, len(page.List))}
	for _, task := range page.List {
		entity, err := taskEntityWithOverdue(l.ctx, base.ProjectManagementTaskLogic, task)
		if err != nil {
			return nil, err
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementTaskPageResponse{Data: data}, nil
}
