package projectmanagementlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
)

type WorkflowUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkflowUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowUpdateLogic {
	return &WorkflowUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowUpdateLogic) WorkflowUpdate(in *borrower.ProjectManagementWorkflowUpdateRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	data, err := workflowDto(in.Data.Name, in.Data.Description, in.Data.Enabled, in.Data.IsDefault, in.Data.Stages, in.Data.Transitions)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	workflow, err := base.ProjectManagementLogic.UpdateWorkflow(l.ctx, actor, in.Data.Id, *data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	entity, err := workflowEntity(workflow)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: entity}, nil
}
