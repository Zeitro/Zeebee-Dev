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

type MemberCandidateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewMemberCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberCandidateLogic {
	return &MemberCandidateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *MemberCandidateLogic) MemberCandidate(in *borrower.ProjectManagementMemberCandidateRequest) (resp *borrower.ProjectManagementMemberCandidateResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	actor, err := requireProjectManagementActor(l.ctx, l.baseLogic)
	if err != nil {
		return nil, err
	}
	members, err := l.baseLogic.ProjectManagementLogic.CandidateProjectMembers(l.ctx, actor, in.Data.ProjectId, in.Data.Keyword)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*dto.ProjectManagementMemberDetailDto, []*borrower.ProjectManagementMemberEntity](&members)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementMemberCandidateResponse{Data: *data}, nil
}
