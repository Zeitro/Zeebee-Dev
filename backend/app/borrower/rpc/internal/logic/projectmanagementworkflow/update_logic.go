package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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

func (l *UpdateLogic) Update(in *borrower.ProjectManagementWorkflowUpdateRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[borrower.ProjectManagementWorkflowUpdateRequestData, dto.ProjectManagementWorkflowUpdateDto](in.Data)
	if err != nil {
		return nil, err
	}
	var workflow *model.ProjectManagementWorkflow
	if in.Data.ProjectId != 0 {
		workflow, err = base.ProjectManagementWorkflowLogic.ProjectUpdate(l.ctx, claim, in.Data.ProjectId, *data)
	} else {
		workflow, err = base.ProjectManagementWorkflowLogic.Update(l.ctx, claim, in.Data.Id, *data)
	}
	if err != nil {
		return nil, err
	}
	entity, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: entity}, nil
}
