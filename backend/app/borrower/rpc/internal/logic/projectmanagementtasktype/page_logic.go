package projectmanagementtasktypelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.ProjectManagementTaskTypePageRequest) (resp *borrower.ProjectManagementTaskTypePageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	input, err := copierx.Trans[borrower.ProjectManagementTaskTypePageRequestData, dto.ProjectManagementTaskTypePageRequestDto](in.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	page, err := l.baseLogic.ProjectManagementTaskTypeLogic.Page(l.ctx, claim.OrganizationId, *input)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data := &borrower.ProjectManagementTaskTypePageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.ProjectManagementTaskTypeEntity, 0, len(page.List))}
	for _, taskType := range page.List {
		entity, err := taskTypeEntity(taskType)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementTaskTypePageResponse{Data: data}, nil
}
