package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type TaskCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCreateLogic {
	return &TaskCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskCreateLogic) TaskCreate(in *borrower.ProjectManagementTaskCreateRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, logic.GetBaseLogic(l.svcCtx))
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
	task, err := logic.GetBaseLogic(l.svcCtx).ProjectManagementLogic.CreateTask(l.ctx, actor, dto.ProjectManagementTaskCreateDto{ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: in.Data.ProjectId, ParentTaskId: in.Data.ParentTaskId, Title: in.Data.Title, Description: in.Data.Description, Priority: in.Data.Priority, StartAt: startAt, EndAt: endAt}, TaskTypeId: in.Data.TaskTypeId, AssigneeIds: in.Data.AssigneeIds})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := taskEntity(task)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
