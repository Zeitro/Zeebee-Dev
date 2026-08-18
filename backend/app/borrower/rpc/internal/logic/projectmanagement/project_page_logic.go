package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectPageLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectPageLogic {
	return &ProjectPageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectPageLogic) ProjectPage(in *borrower.ProjectManagementProjectPageRequest) (resp *borrower.ProjectManagementProjectPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[borrower.ProjectManagementProjectPageRequestData, dto.ProjectManagementProjectPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	page, err := l.baseLogic.ProjectManagementLogic.PageProjects(l.ctx, actor, *data)
	if err != nil {
		return nil, err
	}
	result, err := copierx.Trans[gormx.PageResponse[*model.ProjectManagementProject], borrower.ProjectManagementProjectPageResponseData](&page)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementProjectPageResponse{Data: result}, nil
}
