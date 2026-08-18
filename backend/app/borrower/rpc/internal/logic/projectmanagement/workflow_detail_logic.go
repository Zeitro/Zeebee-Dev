package projectmanagementlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
)

type WorkflowDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowDetailLogic {
	return &WorkflowDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowDetailLogic) WorkflowDetail(in *borrower.ProjectManagementWorkflowDetailRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	workflow, err := base.ProjectManagementLogic.GetWorkflow(l.ctx, actor.OrganizationId, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	entity, err := workflowEntity(workflow)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: entity}, nil
}
