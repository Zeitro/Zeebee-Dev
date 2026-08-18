package project_management

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementMemberCandidateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// list project member candidates
func NewProjectManagementMemberCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementMemberCandidateLogic {
	return &ProjectManagementMemberCandidateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementMemberCandidateLogic) ProjectManagementMemberCandidate(req *types.ProjectManagementMemberCandidateRequest) (resp *types.ProjectManagementMemberCandidateResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMemberCandidateRequest, borrower.ProjectManagementMemberCandidateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.ProjectId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.MemberCandidate(l.ctx, &borrower.ProjectManagementMemberCandidateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementMemberCandidateResponse, types.ProjectManagementMemberCandidateResponse](rpcResp)
}
