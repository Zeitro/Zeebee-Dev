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

func (l *CreateLogic) Create(in *borrower.ProjectManagementTaskCreateRequest) (resp *borrower.ProjectManagementTaskDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	var input dto.ProjectManagementTaskCreateDto
	err = copierx.Assign(in.Data, &input, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	task, err := base.ProjectManagementTaskLogic.Create(l.ctx, claim, input)
	if err != nil {
		return nil, err
	}
	data, err := taskEntityWithOverdue(l.ctx, base.ProjectManagementTaskLogic, task)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementTaskDetailResponse{Data: data}, nil
}
