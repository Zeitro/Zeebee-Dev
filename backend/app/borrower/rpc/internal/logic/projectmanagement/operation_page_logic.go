package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type OperationPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperationPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationPageLogic {
	return &OperationPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperationPageLogic) OperationPage(in *borrower.ProjectManagementOperationPageRequest) (resp *borrower.ProjectManagementOperationPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	filter := dto.ProjectManagementOperationFilterDto{ProjectId: in.Data.ProjectId}
	if in.Data.Filter != nil {
		filter.TaskId = in.Data.Filter.TaskId
		filter.OperationType = in.Data.Filter.OperationType
	}
	page, err := base.ProjectManagementLogic.PageProjectOperations(l.ctx, actor, dto.ProjectManagementOperationPageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit, Filter: filter})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data := &borrower.ProjectManagementOperationPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*borrower.ProjectManagementOperationEntity, 0, len(page.List))}
	for _, operation := range page.List {
		entity, err := operationEntity(operation)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementOperationPageResponse{Data: data}, nil
}
