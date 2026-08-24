package project_management_project

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete project
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.PathIdRequest, borrower.ProjectManagementProjectDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProjectManagementProjectRpcClient.Delete(l.ctx, &borrower.ProjectManagementProjectDeleteRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
