package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type TaskUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskUpdateLogic {
	return &TaskUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskUpdateLogic) TaskUpdate(in *borrower.ProjectManagementTaskUpdateRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	startAt, err := taskTime(in.Data.StartAt)
	if err != nil {
		return nil, err
	}
	endAt, err := taskTime(in.Data.EndAt)
	if err != nil {
		return nil, err
	}
	task, err := base.ProjectManagementLogic.UpdateTask(l.ctx, actor, in.Data.Id, dto.ProjectManagementTaskUpdateDto{Title: in.Data.Title, Description: in.Data.Description, Priority: in.Data.Priority, StartAt: startAt, EndAt: endAt, TaskTypeId: in.Data.TaskTypeId})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := taskEntity(task)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
