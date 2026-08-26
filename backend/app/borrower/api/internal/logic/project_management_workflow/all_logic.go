package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// list enabled project workflows
func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLogic) All(req *types.ProjectManagementWorkflowAllRequest) (resp *types.ProjectManagementWorkflowAllResponse, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.All(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*proto.ProjectManagementWorkflowEntity, []*types.ProjectManagementWorkflowEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowAllResponse{List: *data}, nil

}
