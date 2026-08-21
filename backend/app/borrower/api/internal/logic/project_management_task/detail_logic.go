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

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get project task
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.ProjectManagementTaskDetailRequest) (resp *types.ProjectManagementTaskDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskDetailRequest, borrower.ProjectManagementTaskDetailRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.Detail(l.ctx, &borrower.ProjectManagementTaskDetailRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementTaskEntity, types.ProjectManagementTaskEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTaskDetailResponse{ProjectManagementTaskEntity: *entity}, nil
}
