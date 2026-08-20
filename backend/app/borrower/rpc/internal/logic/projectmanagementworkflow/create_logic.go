package projectmanagementworkflowlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
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

func (l *CreateLogic) Create(in *borrower.ProjectManagementWorkflowCreateRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
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
	workflow, err := base.ProjectManagementWorkflowLogic.Create(l.ctx, claim, *data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	entity, err := workflowEntity(workflow)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementWorkflowDetailResponse{Data: entity}, nil
}
