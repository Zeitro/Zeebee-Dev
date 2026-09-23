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
		if _, exists := bound[stage.Id]; exists {
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
