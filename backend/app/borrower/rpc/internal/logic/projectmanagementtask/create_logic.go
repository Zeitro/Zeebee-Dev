package projectmanagementtasklogic

import (
	"context"
	"strings"
	"time"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func taskTime(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, errorx.New("task time must use RFC3339")
	}
	return &parsed, nil
}

func (l *CreateLogic) Create(in *borrower.ProjectManagementTaskCreateRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
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
	startAt, err := taskTime(in.Data.StartAt)
	if err != nil {
		return nil, err
	}
	endAt, err := taskTime(in.Data.EndAt)
	if err != nil {
		return nil, err
	}
	task, err := base.ProjectManagementTaskLogic.Create(l.ctx, claim, dto.ProjectManagementTaskCreateDto{ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: in.Data.ProjectId, ParentTaskId: in.Data.ParentTaskId, Title: in.Data.Title, Description: in.Data.Description, Priority: in.Data.Priority, StartAt: startAt, EndAt: endAt}, TaskTypeId: in.Data.TaskTypeId, AssigneeIds: in.Data.AssigneeIds})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := taskEntityWithOverdue(l.ctx, base.ProjectManagementTaskLogic, task)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
