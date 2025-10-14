package file

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	// Get borrower id from context
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Call RPC to delete file
	_, err = l.svcCtx.FileRpcClient.DeleteFile(l.ctx, &borrower.FileDeleteRequest{
		Data: &borrower.FileDeleteRequestData{
			Id:         req.Id,
			BorrowerId: borrowerId,
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	return &types.EmptyResponse{}, nil
}
