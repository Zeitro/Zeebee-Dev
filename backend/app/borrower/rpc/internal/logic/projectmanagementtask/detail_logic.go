package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	domainlogic "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(in *borrower.ProjectManagementTaskDetailRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
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
	task, err := base.ProjectManagementTaskLogic.DetailForProject(l.ctx, claim, in.Data.ProjectId, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := taskEntityWithOverdue(l.ctx, logic.GetBaseLogic(l.svcCtx).ProjectManagementTaskLogic, task)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}

func taskEntity(task *model.ProjectManagementTask, overdue *dto.ProjectManagementTaskOverdueDto) (*proto.ProjectManagementTaskEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementTask, proto.ProjectManagementTaskEntity](task,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}
	if task.TaskType != nil && task.TaskType.Id != 0 {
		entity.TaskType, err = copierx.Trans[dto.ProjectManagementTaskTypeSnapshotDto, proto.ProjectManagementTaskTypeSnapshotEntity](task.TaskType)
		if err != nil {
			return nil, err
		}
	}
	if overdue != nil {
		entity.Overdue, entity.OverdueSeconds = overdue.Overdue, overdue.OverdueSeconds
	}
	assignees, err := copierx.Trans[[]*model.ProjectManagementTaskAssignee, []*proto.ProjectManagementTaskAssigneeEntity](&task.Assignees, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	entity.Assignees = *assignees
	return entity, nil
}

func taskEntityWithOverdue(ctx context.Context, taskLogic *domainlogic.ProjectManagementTaskLogic, task *model.ProjectManagementTask) (*proto.ProjectManagementTaskEntity, error) {
	overdue, err := taskLogic.Overdue(ctx, task)
	if err != nil {
		return nil, err
	}
	return taskEntity(task, overdue)
}
