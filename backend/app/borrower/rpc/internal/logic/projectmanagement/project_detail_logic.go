package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectDetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectDetailLogic {
	return &ProjectDetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectDetailLogic) ProjectDetail(in *borrower.ProjectManagementProjectDetailRequest) (resp *borrower.ProjectManagementProjectDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	project, err := l.baseLogic.ProjectManagementLogic.RequireProject(l.ctx, in.Data.Id, actor.OrganizationId)
	if err != nil {
		return nil, err
	}
	if _, err = l.baseLogic.ProjectManagementLogic.RequireMember(l.ctx, project.Id, actor.OrganizationUserId); err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementProject, borrower.ProjectManagementProjectEntity](project)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementProjectDetailResponse{Data: data}, nil
}
