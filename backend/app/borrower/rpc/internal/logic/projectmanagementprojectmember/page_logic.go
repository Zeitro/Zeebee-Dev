package projectmanagementprojectmemberlogic

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
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.ProjectManagementMemberPageRequest) (resp *borrower.ProjectManagementMemberPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := claimValue
	input, err := copierx.Trans[borrower.ProjectManagementMemberPageRequestData, dto.ProjectManagementMemberPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	page, err := l.baseLogic.ProjectManagementMemberLogic.Page(l.ctx, claim, in.Data.ProjectId, *input)
	if err != nil {
		return nil, err
	}
	data := &borrower.ProjectManagementMemberPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.ProjectManagementMemberEntity, 0, len(page.List))}
	for _, member := range page.List {
		entity, err := memberEntity(member)
		if err != nil {
			return nil, err
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementMemberPageResponse{Data: data}, nil
}
