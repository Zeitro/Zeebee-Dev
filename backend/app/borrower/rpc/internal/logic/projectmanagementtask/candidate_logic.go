package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type CandidateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CandidateLogic {
	return &CandidateLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CandidateLogic) Candidate(in *borrower.ProjectManagementTaskCandidateRequest) (*borrower.ProjectManagementTaskCandidateResponse, error) {
	if in.Data == nil {
		return nil, errorx.New("request data is nil")
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	tasks, err := base.ProjectManagementTaskLogic.Candidate(l.ctx, claim, dto.ProjectManagementTaskCandidateDto{EntityType: constants.ProjectManagementOwnerType(in.Data.EntityType), EntityId: in.Data.EntityId, TaskId: in.Data.TaskId, Keyword: in.Data.Keyword})
	if err != nil {
		return nil, err
	}
	data := make([]*borrower.ProjectManagementTaskCandidateEntity, 0, len(tasks))
	for _, task := range tasks {
		data = append(data, &borrower.ProjectManagementTaskCandidateEntity{Id: task.Id, Title: task.Title, Phase: string(task.Phase)})
	}
	return &borrower.ProjectManagementTaskCandidateResponse{Data: data}, nil
}
