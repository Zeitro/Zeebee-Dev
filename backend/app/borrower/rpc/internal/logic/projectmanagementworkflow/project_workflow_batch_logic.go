package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type ProjectWorkflowBatchLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectWorkflowBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowBatchLogic {
	return &ProjectWorkflowBatchLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectWorkflowBatchLogic) ProjectWorkflowBatch(in *borrower.ProjectManagementProjectWorkflowBatchRequest) (resp *borrower.ProjectManagementWorkflowBatchResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	workflows, err := l.baseLogic.ProjectManagementWorkflowLogic.ProjectWorkflowBatch(l.ctx, claim, in.Data.ProjectId, in.Data.SourceWorkflowIds)
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*model.ProjectManagementWorkflow, []*proto.ProjectManagementWorkflowEntity](&workflows, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementWorkflowBatchResponse{List: *list}, nil
}
