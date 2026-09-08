package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
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
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	var input dto.ProjectManagementTaskUpdateDto
	err = copierx.Assign(in.Data, &input, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	task, err := base.ProjectManagementTaskLogic.Update(l.ctx, claim, in.Data.Id, input)
	if err != nil {
		return nil, err
	}
	data, err := taskEntityWithOverdue(l.ctx, logic.GetBaseLogic(l.svcCtx).ProjectManagementTaskLogic, task)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
