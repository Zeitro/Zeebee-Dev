package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
)

type TaskDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskDetailLogic {
	return &TaskDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskDetailLogic) TaskDetail(in *borrower.ProjectManagementTaskDetailRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	task, err := base.ProjectManagementLogic.GetTask(l.ctx, actor, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := taskEntityWithOverdue(l.ctx, logic.GetBaseLogic(l.svcCtx).ProjectManagementLogic, task)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
