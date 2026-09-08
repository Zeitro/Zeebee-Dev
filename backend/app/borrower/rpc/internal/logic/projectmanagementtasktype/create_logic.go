package projectmanagementtasktypelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *borrower.ProjectManagementTaskTypeCreateRequest) (resp *borrower.ProjectManagementTaskTypeDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.ProjectManagementTaskTypeCreateRequestData, dto.ProjectManagementTaskTypeDto](in.Data)
	if err != nil {
		return nil, err
	}
	taskType, err := l.baseLogic.ProjectManagementTaskTypeLogic.Create(l.ctx, claim, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementTaskType, proto.ProjectManagementTaskTypeEntity](taskType, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementTaskTypeDetailResponse{Data: data}, nil
}
