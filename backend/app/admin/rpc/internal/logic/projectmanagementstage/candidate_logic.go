package projectmanagementstagelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

const (
	workflowTypeUnbound = "Unbound"
	workflowTypeBound   = "Bound"
)

type CandidateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CandidateLogic {
	return &CandidateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CandidateLogic) Candidate(in *admin.ProjectManagementStageCandidateRequest) (resp *admin.ProjectManagementStageCandidateResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	if in.Data.WorkflowType != "" && in.Data.WorkflowType != workflowTypeUnbound && in.Data.WorkflowType != workflowTypeBound {
		return nil, errorx.New("workflow type must be Unbound or Bound")
	}
	if in.Data.WorkflowType == workflowTypeBound && in.Data.WorkflowId == 0 {
		return nil, errorx.New("workflow id is required for bound workflow type")
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	stages, err := base.ProjectManagementTemplateLogic.Stages(l.ctx, claim, in.Data.TemplateId)
	if err != nil {
		return nil, err
	}
	bound := make(map[int64]struct{})
	if in.Data.WorkflowId != 0 {
		workflow, err := base.ProjectManagementTemplateLogic.WorkflowDetail(l.ctx, claim, in.Data.TemplateId, in.Data.WorkflowId)
		if err != nil {
			return nil, err
		}
		for _, stage := range workflow.Stages {
			bound[stage.Id] = struct{}{}
		}
	}
	list := make([]*admin.ProjectManagementStageCandidateEntity, 0, len(stages))
	for _, stage := range stages {
		if !includeCandidateStage(bound, stage.Id, in.Data.WorkflowType) {
			continue
		}
		item, err := copierx.Trans[model.ProjectManagementStage, admin.ProjectManagementStageCandidateEntity](stage)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return &admin.ProjectManagementStageCandidateResponse{List: list}, nil
}

func includeCandidateStage(bound map[int64]struct{}, stageID int64, workflowType string) bool {
	_, exists := bound[stageID]
	return exists == (workflowType == workflowTypeBound)
}
