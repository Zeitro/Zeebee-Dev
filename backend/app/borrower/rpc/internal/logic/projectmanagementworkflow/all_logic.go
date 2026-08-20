package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *proto.EmptyRequest) (resp *borrower.ProjectManagementWorkflowAllResponse, err error) {
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	workflows, err := l.baseLogic.ProjectManagementWorkflowLogic.ListEnabled(l.ctx, claim.OrganizationId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	entities := make([]*proto.ProjectManagementWorkflowEntity, 0, len(workflows))
	for _, workflow := range workflows {
		entity, err := workflowEntity(workflow)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	resp = &borrower.ProjectManagementWorkflowAllResponse{Data: entities}

	return resp, nil
}
