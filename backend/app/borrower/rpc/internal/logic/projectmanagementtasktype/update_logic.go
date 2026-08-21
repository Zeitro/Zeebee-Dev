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

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *borrower.ProjectManagementTaskTypeUpdateRequest) (resp *borrower.ProjectManagementTaskTypeDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := claimValue
	input, err := copierx.Trans[borrower.ProjectManagementTaskTypeUpdateRequestData, dto.ProjectManagementTaskTypeDto](in.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	taskType, err := l.baseLogic.ProjectManagementTaskTypeLogic.Update(l.ctx, claim, in.Data.Id, *input)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := copierx.Trans[model.ProjectManagementTaskType, proto.ProjectManagementTaskTypeEntity](taskType, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskTypeDetailResponse{Data: data}, nil
}
