package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/common/errorx"
)

type ProjectWorkflowUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProjectWorkflowUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowUpdateLogic {
	return &ProjectWorkflowUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectWorkflowUpdateLogic) ProjectWorkflowUpdate(in *borrower.ProjectManagementProjectWorkflowUpdateRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	data, err := workflowDto(in.Data.Name, in.Data.Description, true, false, in.Data.Stages, in.Data.Transitions)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	workflow, err := base.ProjectManagementLogic.UpdateProjectWorkflow(l.ctx, actor, in.Data.ProjectId, *data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	entity, err := workflowEntity(workflow)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: entity}, nil
}
