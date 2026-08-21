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

type MyPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyPageLogic {
	return &MyPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyPageLogic) MyPage(in *borrower.ProjectManagementMyTaskPageRequest) (resp *borrower.ProjectManagementTaskPageResponse, err error) {
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
	filter := dto.ProjectManagementTaskFilterDto{}
	if in.Data.Filter != nil {
		filter.ProjectId = in.Data.Filter.ProjectId
		filter.StageId = in.Data.Filter.StageId
		filter.TaskTypeId = in.Data.Filter.TaskTypeId
		filter.AssigneeId = in.Data.Filter.AssigneeId
		filter.Keyword = in.Data.Filter.Keyword
		filter.Overdue = in.Data.Filter.Overdue
	}
	page, err := base.ProjectManagementTaskLogic.MyPage(l.ctx, claim, dto.ProjectManagementMyTaskPageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit, Filter: filter})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data := &borrower.ProjectManagementTaskPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.ProjectManagementTaskEntity, 0, len(page.List))}
	for _, task := range page.List {
		entity, err := taskEntityWithOverdue(l.ctx, base.ProjectManagementTaskLogic, task)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementTaskPageResponse{Data: data}, nil
}
