package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type ViewLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewViewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ViewLogic {
	return &ViewLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ViewLogic) View(in *borrower.ProjectManagementTaskViewRequest) (resp *proto.ProjectManagementTaskViewEntity, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	view, err := l.baseLogic.ProjectManagementTaskLogic.View(l.ctx, claim, in.Data.Id)
	if err != nil {
		return nil, err
	}
	result, err := copierx.Trans[dto.ProjectManagementTaskViewDto, proto.ProjectManagementTaskViewEntity](view, copierx.SqlTimePtrToRFC3339Str())
	if err != nil {
		return nil, err
	}
	result.Tasks = make([]*proto.ProjectManagementTaskEntity, 0, len(view.Tasks))
	for _, task := range view.Tasks {
		entity, convertErr := taskEntity(task, nil)
		if convertErr != nil {
			return nil, convertErr
		}
		result.Tasks = append(result.Tasks, entity)
	}
	return result, nil
}
