package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

type WorkflowAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewWorkflowAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowAllLogic {
	return &WorkflowAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *WorkflowAllLogic) WorkflowAll(in *borrower.ProjectManagementWorkflowAllRequest) (resp *borrower.ProjectManagementWorkflowAllResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	workflows, err := l.baseLogic.ProjectManagementLogic.ListEnabledWorkflows(l.ctx, actor.OrganizationId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	data, err := copierx.Trans[[]*model.ProjectManagementWorkflow, []*borrower.ProjectManagementWorkflowEntity](&workflows)
	if err != nil {
		return nil, err
	}
	resp = &borrower.ProjectManagementWorkflowAllResponse{Data: *data}

	return resp, nil
}
