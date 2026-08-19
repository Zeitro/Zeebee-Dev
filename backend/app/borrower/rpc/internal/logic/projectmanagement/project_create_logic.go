package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type ProjectCreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectCreateLogic {
	return &ProjectCreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectCreateLogic) ProjectCreate(in *borrower.ProjectManagementProjectCreateRequest) (resp *borrower.ProjectManagementProjectDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	projectData, err := copierx.Trans[borrower.ProjectManagementProjectCreateRequestData, dto.ProjectManagementProjectDto](in.Data)
	if err != nil {
		return nil, err
	}
	project, err := l.baseLogic.ProjectManagementLogic.CreateProject(l.ctx, actor, dto.ProjectManagementProjectCreateDto{ProjectManagementProjectDto: *projectData, WorkflowTemplateId: in.Data.WorkflowTemplateId})
	if err != nil {
		return nil, err
	}
	data, err := projectEntity(project)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementProjectDetailResponse{Data: data}, nil
}
