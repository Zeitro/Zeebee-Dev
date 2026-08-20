package projectmanagementworkflowlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/model"

	"zeebee-backend/common/errorx"
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

func (l *UpdateLogic) Update(in *borrower.ProjectManagementWorkflowUpdateRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	data, err := workflowDto(in.Data.Name, in.Data.Description, in.Data.Enabled, in.Data.IsDefault, in.Data.Stages, in.Data.Transitions)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	var workflow *model.ProjectManagementWorkflow
	if in.Data.ProjectId != 0 {
		workflow, err = base.ProjectManagementWorkflowLogic.ProjectUpdate(l.ctx, claim, in.Data.ProjectId, *data)
	} else {
		workflow, err = base.ProjectManagementWorkflowLogic.Update(l.ctx, claim, in.Data.Id, *data)
	}
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	entity, err := workflowEntity(workflow)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: entity}, nil
}
