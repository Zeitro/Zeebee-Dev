package projectmanagementlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
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
	result, err := copierx.Trans[gormx.PageResponse[*model.ProjectManagementWorkflow], borrower.ProjectManagementWorkflowPageResponseData](&page)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementWorkflowPageResponse{Data: result}, nil
}
