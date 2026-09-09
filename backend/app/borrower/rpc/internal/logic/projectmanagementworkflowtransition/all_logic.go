package projectmanagementworkflowtransitionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLogic) All(in *borrower.ProjectManagementWorkflowTransitionAllRequest) (resp *borrower.ProjectManagementWorkflowTransitionAllResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	transitions, err := base.ProjectManagementWorkflowLogic.TransitionAll(l.ctx, claim.OrganizationId, in.Data.WorkflowId, in.Data.FromStageId)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*model.ProjectManagementWorkflowTransition, []*proto.ProjectManagementWorkflowTransitionEntity](&transitions, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementWorkflowTransitionAllResponse{Data: *data}, nil
}
