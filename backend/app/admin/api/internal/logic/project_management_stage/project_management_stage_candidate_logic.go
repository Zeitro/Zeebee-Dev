package project_management_stage

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementStageCandidateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// stages candidate
func NewProjectManagementStageCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementStageCandidateLogic {
	return &ProjectManagementStageCandidateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementStageCandidateLogic) ProjectManagementStageCandidate(req *types.ProjectManagementStageCandidateRequest) (resp *types.ProjectManagementStageCandidateResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementStageCandidateRequest, admin.ProjectManagementStageCandidateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementStageRpcClient.Candidate(l.ctx, &admin.ProjectManagementStageCandidateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*admin.ProjectManagementStageCandidateEntity, []types.ProjectManagementStageCandidateEntity](&rpcResp.List)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementStageCandidateResponse{List: *list}, nil
}
