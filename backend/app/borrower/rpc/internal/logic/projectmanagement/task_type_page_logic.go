package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type TaskTypePageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewTaskTypePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskTypePageLogic {
	return &TaskTypePageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *TaskTypePageLogic) TaskTypePage(in *borrower.ProjectManagementTaskTypePageRequest) (resp *borrower.ProjectManagementTaskTypePageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.ProjectManagementTaskTypePageRequestData, dto.ProjectManagementTaskTypePageRequestDto](in.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	page, err := l.baseLogic.ProjectManagementLogic.PageTaskTypes(l.ctx, actor.OrganizationId, *input)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := copierx.Trans[gormx.PageResponse[*model.ProjectManagementTaskType], borrower.ProjectManagementTaskTypePageResponseData](&page)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskTypePageResponse{Data: data}, nil
}
