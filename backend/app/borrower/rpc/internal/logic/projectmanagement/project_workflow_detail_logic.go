package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
)

type ProjectWorkflowDetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowDetailLogic {
	return &ProjectWorkflowDetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectWorkflowDetailLogic) ProjectWorkflowDetail(in *borrower.ProjectManagementProjectWorkflowDetailRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	workflow, err := l.baseLogic.ProjectManagementLogic.GetProjectWorkflow(l.ctx, actor, in.Data.ProjectId)
	if err != nil {
		return nil, err
	}
	data, err := workflowEntity(workflow)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: data}, nil
}
