package projectmanagementworkflowstagelogic

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

func (l *AllLogic) All(in *borrower.ProjectManagementWorkflowStageAllRequest) (resp *borrower.ProjectManagementWorkflowStageAllResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	stages, err := base.ProjectManagementWorkflowLogic.StageAll(l.ctx, claim.OrganizationId, in.Data.WorkflowId)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*model.ProjectManagementWorkflowStage, []*proto.ProjectManagementWorkflowStageEntity](&stages, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementWorkflowStageAllResponse{Data: *data}, nil
}
