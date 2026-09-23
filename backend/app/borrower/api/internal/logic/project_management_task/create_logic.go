package project_management_task

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

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create project task
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProjectManagementTaskCreateRequest) (resp *types.ProjectManagementTaskDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskCreateRequest, borrower.ProjectManagementTaskCreateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.Create(l.ctx, &borrower.ProjectManagementTaskCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskDetailResponse{ProjectManagementTaskEntity: *entity}, nil
}
