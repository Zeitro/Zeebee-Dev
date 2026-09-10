package project_management_project

import (
	"context"

	workflowlogic "zeebee-backend/app/borrower/api/internal/logic/project_management_workflow"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkflowUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// workflow update
func NewWorkflowUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowUpdateLogic {
	return &WorkflowUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowUpdateLogic) WorkflowUpdate(req *types.ProjectManagementProjectWorkflowUpdateRequest) (resp *types.ProjectManagementWorkflowDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowUpdateRequest, borrower.ProjectManagementWorkflowUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	data.ProjectId = req.Id
	return workflowlogic.UpdateWorkflow(l.ctx, l.svcCtx, data)
}
