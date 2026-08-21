package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
)

type TransitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransitionLogic {
	return &TransitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransitionLogic) Transition(in *borrower.ProjectManagementTaskTransitionRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := claimValue
	if _, err = base.ProjectManagementTaskLogic.DetailForProject(l.ctx, claim, in.Data.ProjectId, in.Data.TaskId); err != nil {
		return nil, err
	}
	task, err := base.ProjectManagementTaskLogic.Transition(l.ctx, claim, in.Data.TaskId, in.Data.TargetStageId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := taskEntityWithOverdue(l.ctx, logic.GetBaseLogic(l.svcCtx).ProjectManagementTaskLogic, task)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
