package project_management

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementTaskTypeDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete project task type
func NewProjectManagementTaskTypeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskTypeDeleteLogic {
	return &ProjectManagementTaskTypeDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskTypeDeleteLogic) ProjectManagementTaskTypeDelete(req *types.ProjectManagementTaskTypeDeleteRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskTypeDeleteRequest, borrower.ProjectManagementTaskTypeDeleteRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if _, err = l.svcCtx.ProjectManagementRpcClient.TaskTypeDelete(l.ctx, &borrower.ProjectManagementTaskTypeDeleteRequest{Data: data}); err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
