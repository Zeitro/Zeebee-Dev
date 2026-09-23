package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type ProjectWorkflowStageUpdateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectWorkflowStageUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowStageUpdateLogic {
	return &ProjectWorkflowStageUpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectWorkflowStageUpdateLogic) ProjectWorkflowStageUpdate(in *borrower.ProjectManagementProjectWorkflowStageUpdateRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	stages, err := copierx.Trans[[]*borrower.ProjectManagementWorkflowUpdateStage, []dto.ProjectManagementWorkflowUpdateStageDto](&in.Data.Stages)
	if err != nil {
		return nil, err
	}
	workflow, err := l.baseLogic.ProjectManagementWorkflowLogic.ProjectWorkflowStageUpdate(l.ctx, claim, in.Data.ProjectId, in.Data.WorkflowId, dto.ProjectManagementWorkflowUpdateDto{Stages: *stages})
	if err != nil {
		return nil, err
	}
	entity, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: entity}, nil
}
