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

type ProjectManagementProjectDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete project
func NewProjectManagementProjectDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectDeleteLogic {
	return &ProjectManagementProjectDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectDeleteLogic) ProjectManagementProjectDelete(req *types.ProjectManagementProjectDeleteRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectDeleteRequest, borrower.ProjectManagementProjectDeleteRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	_, err = l.svcCtx.ProjectManagementRpcClient.ProjectDelete(l.ctx, &borrower.ProjectManagementProjectDeleteRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
