package project_management_task_type

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.PathIdRequest) (resp *types.ProjectManagementTaskTypeDetailResponse, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementTaskTypeRpcClient.Detail(l.ctx, &proto.PathIdRequest{Data: &proto.PathIdRequestData{Id: req.Id}})
	if err != nil {
		return nil, err
	}
	entity, err := copierx.Trans[proto.ProjectManagementTaskTypeEntity, types.ProjectManagementTaskTypeEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskTypeDetailResponse{ProjectManagementTaskTypeEntity: *entity}, nil
}
