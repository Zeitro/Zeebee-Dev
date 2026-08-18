package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type TaskTypeCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewTaskTypeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskTypeCreateLogic {
	return &TaskTypeCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *TaskTypeCreateLogic) TaskTypeCreate(in *borrower.ProjectManagementTaskTypeCreateRequest) (resp *borrower.ProjectManagementTaskTypeDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.ProjectManagementTaskTypeCreateRequestData, dto.ProjectManagementTaskTypeDto](in.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	taskType, err := l.baseLogic.ProjectManagementLogic.CreateTaskType(l.ctx, actor, *input)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := copierx.Trans[model.ProjectManagementTaskType, borrower.ProjectManagementTaskTypeEntity](taskType)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskTypeDetailResponse{Data: data}, nil
}
