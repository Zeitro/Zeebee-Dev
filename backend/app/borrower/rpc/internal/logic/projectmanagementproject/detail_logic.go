package projectmanagementprojectlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

type DetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *proto.PathIdRequest) (resp *borrower.ProjectManagementProjectDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	project, err := l.baseLogic.ProjectManagementLogic.RequireProject(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	if _, err = l.baseLogic.ProjectManagementLogic.RequireMember(l.ctx, project.Id, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementProject, proto.ProjectManagementProjectEntity](project, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	workflow, err := l.baseLogic.ProjectManagementWorkflowLogic.ProjectDetail(l.ctx, claim, project.Id)
	if err != nil {
		return nil, err
	}
	data.Workflow, err = copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementProjectDetailResponse{Data: data}, nil
}
