package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update project workflow
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProjectManagementWorkflowUpdateRequest) (resp *types.ProjectManagementWorkflowDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowUpdateRequest, borrower.ProjectManagementWorkflowUpdateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Update(l.ctx, &borrower.ProjectManagementWorkflowUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowDetailResponse{ProjectManagementWorkflowEntity: *entity}, nil
}
