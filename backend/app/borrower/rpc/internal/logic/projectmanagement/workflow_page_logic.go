package projectmanagementlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type WorkflowPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkflowPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowPageLogic {
	return &WorkflowPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowPageLogic) WorkflowPage(in *borrower.ProjectManagementWorkflowPageRequest) (resp *borrower.ProjectManagementWorkflowPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[borrower.ProjectManagementWorkflowPageRequestData, dto.ProjectManagementWorkflowPageRequestDto](in.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	page, err := base.ProjectManagementLogic.PageWorkflows(l.ctx, actor.OrganizationId, *data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	result := &borrower.ProjectManagementWorkflowPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*borrower.ProjectManagementWorkflowPageEntity, 0, len(page.List))}
	for _, workflow := range page.List {
		entity, err := workflowPageEntity(workflow)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		result.List = append(result.List, entity)
	}
	return &borrower.ProjectManagementWorkflowPageResponse{Data: result}, nil
}
