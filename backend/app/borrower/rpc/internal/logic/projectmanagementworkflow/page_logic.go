package projectmanagementworkflowlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(in *borrower.ProjectManagementWorkflowPageRequest) (resp *borrower.ProjectManagementWorkflowPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	if in.Data.Filter == nil {
		return nil, errorx.New("request filter is nil")
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	page, err := base.ProjectManagementWorkflowLogic.Page(l.ctx, claim.OrganizationId, dto.ProjectManagementWorkflowPageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit, Filter: dto.ProjectManagementWorkflowFilterDto{Keyword: in.Data.Filter.Keyword, Enabled: in.Data.Filter.Enabled}})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	result := &borrower.ProjectManagementWorkflowPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.ProjectManagementWorkflowPageEntity, 0, len(page.List))}
	for _, workflow := range page.List {
		entity, err := workflowPageEntity(workflow)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		result.List = append(result.List, entity)
	}
	return &borrower.ProjectManagementWorkflowPageResponse{Data: result}, nil
}
