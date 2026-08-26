package project_management_project_member

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CandidateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// list project member candidates
func NewCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CandidateLogic {
	return &CandidateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CandidateLogic) Candidate(req *types.ProjectManagementMemberCandidateRequest) (resp *types.ProjectManagementMemberCandidateResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMemberCandidateRequest, borrower.ProjectManagementMemberCandidateRequestData](req)
	if err != nil {
		return nil, err
	}
	data.ProjectId = req.ProjectId
	rpcResp, err := l.svcCtx.ProjectManagementProjectMemberRpcClient.Candidate(l.ctx, &borrower.ProjectManagementMemberCandidateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	listData, err := copierx.Trans[[]*borrower.ProjectManagementMemberCandidateEntity, []*types.ProjectManagementMemberCandidateEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementMemberCandidateResponse{List: *listData}, nil
}
