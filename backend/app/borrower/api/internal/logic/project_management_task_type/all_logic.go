package project_management_task_type

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

// all
func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLogic) All(req *types.ProjectManagementTaskTypeAllRequest) (resp *types.ProjectManagementTaskTypeAllResponse, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementTaskTypeRpcClient.All(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*proto.ProjectManagementTaskTypeEntity, []*types.ProjectManagementTaskTypeEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskTypeAllResponse{List: *data}, nil
}
