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

type ProjectManagementTaskDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete project task
func NewProjectManagementTaskDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTaskDeleteLogic {
	return &ProjectManagementTaskDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTaskDeleteLogic) ProjectManagementTaskDelete(req *types.ProjectManagementTaskDeleteRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTaskDeleteRequest, borrower.ProjectManagementTaskDeleteRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	_, err = l.svcCtx.ProjectManagementRpcClient.TaskDelete(l.ctx, &borrower.ProjectManagementTaskDeleteRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
