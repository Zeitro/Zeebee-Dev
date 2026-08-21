package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(in *borrower.ProjectManagementTaskUpdateRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
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
	if _, err = base.ProjectManagementTaskLogic.DetailForProject(l.ctx, claim, in.Data.ProjectId, in.Data.Id); err != nil {
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
	task, err := base.ProjectManagementTaskLogic.Update(l.ctx, claim, in.Data.Id, dto.ProjectManagementTaskUpdateDto{Title: in.Data.Title, Description: in.Data.Description, Priority: in.Data.Priority, StartAt: startAt, EndAt: endAt, TaskTypeId: in.Data.TaskTypeId})
	if err != nil {
		return nil, err
	}
	data, err := taskEntityWithOverdue(l.ctx, logic.GetBaseLogic(l.svcCtx).ProjectManagementTaskLogic, task)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
