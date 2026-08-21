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

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get project workflow
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.PathIdRequest) (resp *types.ProjectManagementWorkflowDetailResponse, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Detail(l.ctx, &borrower.ProjectManagementWorkflowDetailRequest{Data: &borrower.ProjectManagementWorkflowDetailRequestData{Id: req.Id}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowDetailResponse{ProjectManagementWorkflowEntity: *entity}, nil
}
