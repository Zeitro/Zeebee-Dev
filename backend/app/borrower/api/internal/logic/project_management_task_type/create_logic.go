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

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create project task type
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProjectManagementTaskTypeCreateRequest) (resp *types.ProjectManagementTaskTypeDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTypeCreateRequest, borrower.ProjectManagementTaskTypeCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementTaskTypeRpcClient.Create(l.ctx, &borrower.ProjectManagementTaskTypeCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementTaskTypeEntity, types.ProjectManagementTaskTypeEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskTypeDetailResponse{ProjectManagementTaskTypeEntity: *entity}, nil
}
