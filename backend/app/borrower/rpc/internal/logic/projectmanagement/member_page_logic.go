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

type MemberPageLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewMemberPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberPageLogic {
	return &MemberPageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *MemberPageLogic) MemberPage(in *borrower.ProjectManagementMemberPageRequest) (resp *borrower.ProjectManagementMemberPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.ProjectManagementMemberPageRequestData, dto.ProjectManagementMemberPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	page, err := l.baseLogic.ProjectManagementLogic.PageProjectMembers(l.ctx, actor, in.Data.ProjectId, *input)
	if err != nil {
		return nil, err
	}
	data := &borrower.ProjectManagementMemberPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*borrower.ProjectManagementMemberEntity, 0, len(page.List))}
	for _, member := range page.List {
		entity, err := memberEntity(member)
		if err != nil {
			return nil, err
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementMemberPageResponse{Data: data}, nil
}
