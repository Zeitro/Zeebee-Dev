package project_management_task

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

func NewCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CandidateLogic {
	return &CandidateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CandidateLogic) Candidate(req *types.ProjectManagementTaskCandidateRequest) (*types.ProjectManagementTaskCandidateResponse, error) {
	data, err := copierx.Trans[types.ProjectManagementTaskCandidateRequest, borrower.ProjectManagementTaskCandidateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.Candidate(l.ctx, &borrower.ProjectManagementTaskCandidateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	list, err := copierx.Trans[[]*borrower.ProjectManagementTaskCandidateEntity, []*types.ProjectManagementTaskCandidateEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskCandidateResponse{List: *list}, nil
}
