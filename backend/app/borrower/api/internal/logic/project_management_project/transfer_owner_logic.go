package project_management_project

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferOwnerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// transfer owner
func NewTransferOwnerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferOwnerLogic {
	return &TransferOwnerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferOwnerLogic) TransferOwner(req *types.ProjectManagementProjectOwnerTransferRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.ProjectManagementProjectRpcClient.TransferOwner(l.ctx, &borrower.ProjectManagementProjectOwnerTransferRequest{Data: &borrower.ProjectManagementProjectOwnerTransferRequestData{
		ProjectId:          req.ProjectId,
		OrganizationUserId: req.OrganizationUserId,
	}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
