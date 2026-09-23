package projectmanagementworkflowlogic

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

type ProjectWorkflowPageLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectWorkflowPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowPageLogic {
	return &ProjectWorkflowPageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectWorkflowPageLogic) ProjectWorkflowPage(in *borrower.ProjectManagementProjectWorkflowPageRequest) (resp *borrower.ProjectManagementWorkflowPageResponse, err error) {
	if in.Data == nil || in.Data.Filter == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	page, err := l.baseLogic.ProjectManagementWorkflowLogic.ProjectWorkflowPage(l.ctx, claim, in.Data.Filter.ProjectId, in.Data.Page, in.Data.Limit)
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*model.ProjectManagementWorkflow, []*proto.ProjectManagementWorkflowEntity](&page.List, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementWorkflowPageResponse{Total: page.Total, Page: page.Page, Limit: page.Limit, TotalPage: page.TotalPage, List: *list}, nil
}
