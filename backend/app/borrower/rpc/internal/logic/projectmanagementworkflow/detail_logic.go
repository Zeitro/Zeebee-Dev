package projectmanagementworkflowlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(in *borrower.ProjectManagementWorkflowDetailRequest) (resp *borrower.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := claimValue
	var workflow *model.ProjectManagementWorkflow
	if in.Data.ProjectId != 0 {
		workflow, err = base.ProjectManagementWorkflowLogic.ProjectDetail(l.ctx, claim, in.Data.ProjectId)
	} else {
		workflow, err = base.ProjectManagementWorkflowLogic.Detail(l.ctx, claim.OrganizationId, in.Data.Id)
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
