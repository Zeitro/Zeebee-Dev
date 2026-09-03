package projectmanagementprojectlogic

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
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *borrower.ProjectManagementProjectUpdateRequest) (resp *borrower.ProjectManagementProjectDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.ProjectManagementProjectUpdateRequestData, dto.ProjectManagementProjectUpdateDto](in.Data)
	if err != nil {
		return nil, err
	}
	project, err := l.baseLogic.ProjectManagementProjectLogic.Update(l.ctx, claim, in.Data.Id, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementProject, proto.ProjectManagementProjectEntity](project, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementProjectDetailResponse{Data: data}, nil
}
