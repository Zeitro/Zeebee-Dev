package projectmanagementprojectmemberlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type CandidateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CandidateLogic {
	return &CandidateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CandidateLogic) Candidate(in *borrower.ProjectManagementMemberCandidateRequest) (resp *borrower.ProjectManagementMemberCandidateResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	members, err := l.baseLogic.ProjectManagementMemberLogic.Candidate(l.ctx, claim, in.Data.ProjectId, in.Data.Keyword, in.Data.NeedAll)
	if err != nil {
		return nil, err
	}
	data := make([]*proto.ProjectManagementMemberEntity, 0, len(members))
	for _, member := range members {
		entity, err := memberEntity(member)
		if err != nil {
			return nil, err
		}
		data = append(data, entity)
	}
	return &borrower.ProjectManagementMemberCandidateResponse{Data: data}, nil
}
