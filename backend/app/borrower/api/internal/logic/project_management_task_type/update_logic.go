package project_management_task_type

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

// update project task type
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProjectManagementTaskTypeUpdateRequest) (resp *types.ProjectManagementTaskTypeDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTypeUpdateRequest, borrower.ProjectManagementTaskTypeUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTaskTypeRpcClient.Update(l.ctx, &borrower.ProjectManagementTaskTypeUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementTaskTypeEntity, types.ProjectManagementTaskTypeEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskTypeDetailResponse{ProjectManagementTaskTypeEntity: *entity}, nil
}
